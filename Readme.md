# CS670 Assignment 1 — Secure User-Profile Update (Additive MPC)

This repository contains a runnable reference implementation for the assignment.  
It demonstrates **additive secret sharing** over the ring `Z_(2^64)` together with a **semi-honest dealer (P2)** that distributes Beaver triples for secure multiplications.

---

## How it Works (Short Report)

### Secret Sharing
- A secret value `x` is represented as two additive shares:
  ```
  x = x0 + x1   (mod 2^64)
  ```
  - Party **P0** holds `x0`
  - Party **P1** holds `x1`
- Vectors are shared element-by-element.  
- The dealer **P2** generates plaintext matrices:
  - `U ∈ Z^(m×k)` = user feature matrix
  - `V ∈ Z^(n×k)` = item feature matrix  
  Each entry is sampled randomly between 1 and 7 (small values make debugging easier).  
- P2 splits each value into shares and distributes them to P0 and P1.

---

### Secure Inner Product
We want to compute:
```
s = <u_i, v_j> = Σ_t (u_i[t] * v_j[t])
```
without revealing `u_i` or `v_j`.

This is done using **Beaver triples**:
- A triple `(a, b, c)` with `c = a * b` is pre-shared between P0 and P1.  
- For inputs `(x, y)`:
  1. Each party computes differences:
     ```
     e = x - a
     f = y - b
     ```
  2. Parties **open** (exchange) `e` and `f`.  
  3. Each party computes its share of `z = x * y` as:
     ```
     z_share = c_share + e * b_share + f * a_share + (if P0) e * f
     ```
- Repeating this for all coordinates and summing the results yields shares of `s`.

---

### Secure Update
After computing `s`, we update the user vector `u_i`:
```
delta = 1 - s
u_i <- u_i + v_j * delta
```

- The constant `1` is public, so it is shared as `(1, 0)`.  
- To compute `v_j * delta`, we use a **vector–scalar Beaver triple**:
  ```
  (a_vec, b, c_vec)   with   c_vec = a_vec * b
  ```
- Each party computes:
  ```
  e_vec = v_j - a_vec
  f = delta - b
  ```
  and opens `(e_vec, f)`.  
- Output shares are:
  ```
  z_share = c_vec_share + f * a_vec_share + e_vec * b_share + (if P0) f * e_vec
  ```
- Finally, each party adds `z_share` into its share of `u_i`.

---

### Communication & Efficiency
- Each secure multiplication (scalar or vector–scalar) requires **one round** of communication (opening `e` and `f`).  
- The inner product needs:
  - `k` scalar triples (one per coordinate).  
- The update needs:
  - One vector–scalar triple.  
- In practice:
  - Multiple `e,f` values can be **batched** to reduce communication.  
  - Beaver triples can be pre-generated offline for better efficiency.  
  - Extensions include replicated secret sharing for stronger security.

---

## Build & Run

**Requirements:**  
- Docker  
- docker-compose  

### Build
```bash
docker-compose build
```

### Run
```bash
docker-compose up
```

By default, the compose file runs with:
```
m = 4, n = 4, k = 5, q = 1
```
This computes a single update for `(i=0, j=0)`.  

Check logs of the container `a1_p0` for the reconstructed values.

---

### Custom Sizes
To change the problem dimensions, edit the `p2` service command in `docker-compose.yml` and adjust the `party` commands accordingly.  

- Ports are fixed as follows:
  - P0 connects to `p2:9002`
  - P1 connects to `p2:9003`
  - Parties connect to each other on `9001`

---

## Files

- **`shares.hpp`**  
  Ring operations, vector shares, Beaver triple structures.  

- **`protocol.hpp`**  
  Socket helpers, open/aggregate helpers, and MPC protocol kernels.  

- **`p2.cpp`**  
  Dealer: samples plaintext `U, V`, splits them into shares, and distributes Beaver triples.  

- **`party.cpp`**  
  Logic for P0/P1 to receive shares, run secure dot-product and update, and print shares/reconstructed values.  

- **`Dockerfile`, `docker-compose.yml`**  
  Containerization and orchestration scripts.

---

## Debug Output

- **Dealer (P2)**:
  - Prints the plaintext matrices `U` and `V`.  
  - Prints all generated Beaver triples (dot-product and vector–scalar).  

- **Parties (P0, P1)**:
  - Print their shares of the updated `u_i`.  

- **Party P0**:
  - Prints the reconstructed inner product `<u_i, v_j>`  
  - Prints the reconstructed updated user vector `u_i`.  

This allows you to verify correctness while still keeping the protocol secure.

---

## Notes

- **Arithmetic domain**:  
  All operations are done in `Z_(2^64)` (unsigned 64-bit wraparound).  
  Multiplications use 128-bit temporaries before reduction to avoid overflow.  

- **Security model**:  
  Semi-honest adversary model. Parties follow the protocol but may try to infer extra info.  
  Only masked differences `e,f` are revealed during secure multiplications.  

- **Networking**:  
  Simple blocking TCP sockets using Boost.Asio.  

- **Extensibility**:  
  - Easy to modify for multiple queries or different `(i,j)` selection.  
  - Can be extended with offline triple generation or replicated secret sharing.  
  - Serves as a baseline reference implementation for further research or assignments.

---
