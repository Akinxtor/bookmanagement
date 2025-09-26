# CS670 Assignment 1 — Secure User-Profile Update (Additive MPC)

This repo provides a runnable reference implementation for the assignment spec. It uses additive secret sharing over the ring \( \mathbb{Z}_{2^{64}} \) and a semi-honest dealer \(P2\) to distribute Beaver triples for secure multiplications.

## How it works (short report)

### Secret sharing
We represent a secret \(x\) as a sum of two shares \(x = x_0 + x_1 \pmod{2^{64}}\). Party \(P0\) holds \(x_0\) and \(P1\) holds \(x_1\). Vectors are shared component-wise. The dealer \(P2\) samples the plaintext matrices \(U \in \mathbb{Z}^{m\times k}\) and \(V \in \mathbb{Z}^{n\times k}\) with small entries (1–7 for readability) and splits them into additive shares sent to the two parties.

### Secure inner product
To compute \( s=\langle u_i, v_j\rangle = \sum_t u_{it} v_{jt}\) on shares, we use Beaver triples. For each coordinate product we have a triple \((a,b,c)\) with \(c=a\cdot b\) secret-shared between \(P0\) and \(P1\).
Each party computes \(e = x-a\) and \(f = y-b\) locally on shares and **opens** only \(e,f\) to the peer. Then the output share is
\[ z_i = c_i + e b_i + f a_i + [i=0]\cdot e f. \]
Summing the per-coordinate results yields a share of \(s\).

### Secure update
We set \( \delta = 1 - s\), where 1 is public so it is shared as \((1,0)\). To compute \(u_i \leftarrow u_i + v_j\cdot \delta\) we use a **vector–scalar** Beaver triple \((\mathbf{a}, b, \mathbf{c}=\mathbf{a}\cdot b)\), opening \( \mathbf{e}=\mathbf{v}-\mathbf{a}\) and \(f=\delta-b\), and producing shares of \(\mathbf{v}\cdot\delta\) component-wise using the same formula. Finally we add the result into the current share of \(u_i\).

### Communication rounds & efficiency
Each secure scalar or vector–scalar multiplication takes one round to open \(e,f\) (two small messages per direction). The inner product needs \(k\) scalar triples; the update needs one vector–scalar triple. In practice these opens can be **batched** across coordinates; our implementation already batches the vector open for \( \mathbf{e}\). Extensions include pre-generating triples offline and supporting replicated secret sharing for the bonus part.

## Build & Run

Requirements: Docker and docker-compose.

```bash
docker-compose build
docker-compose up
```

By default, the compose file runs with `m=4, n=4, k=5, q=1` and computes a single update for `(i=0, j=0)`.
Look at the logs of `a1_p0` for the reconstructed values (for debugging only, as requested).

### Custom sizes
You can override the arguments by editing the `p2` service command in `docker-compose.yml`, and adjusting the `party` commands accordingly (ports are fixed: P0 connects to `p2:9002`, P1 to `p2:9003`, and parties connect to each other on `9001`).

## Files

- `shares.hpp`: ring operations, vector shares, and Beaver triple types.
- `protocol.hpp`: sync socket helpers, open/aggregate helpers, and MPC kernels.
- `p2.cpp`: dealer that samples plaintext `U,V`, splits shares, and distributes Beaver triples.
- `party.cpp`: logic for P0/P1 to receive shares, run MPC for dot-product and update, print shares and reconstructed values.
- `Dockerfile`, `docker-compose.yml`: containerization per the spec.

## Debug Output

- The Dealer `P2` prints the shares of scalar and vector beaver triples along with initial U and V vector.
- Each party prints its share of the updated `u_i`.
- Party `P0` also prints the reconstructed inner product `<u_i, v_j>` and the reconstructed `u_i` **after** the update, as required for debugging.

## Notes

- Arithmetic is in \( \mathbb{Z}_{2^{64}} \) (wrap-around). This matches the assignment’s integer-domain requirement and avoids overflows by using 128-bit temporaries in `mod_mul`.
- The code is written to be clear and easy to modify for multiple queries or different `(i,j)` selection; wiring in a real query source or a client is straightforward.
- Security model: semi-honest parties; only masked values \(e,f\) are opened.
- Networking: simple blocking I/O over TCP using Boost.Asio.
