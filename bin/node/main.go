package main

// **Tech Stack**: Go (crypto/sha256, net, fmt)
// - **Steps**:
//     1. Create a blockchain structure where each block represents a vote.
//     2. Use `crypto/sha256` to hash block data and create cryptographic links between blocks.
//     3. Implement basic proof-of-work by making nodes solve a computational problem before adding new blocks.
//     4. Use `net` to implement peer-to-peer communication, allowing nodes to broadcast new blocks to each other.
//     5. Build a simple CLI to interact with the blockchain (e.g., cast votes, view results).
// - **Enhancements**:
//     - Implement a consensus algorithm to resolve chain conflicts and ensure tamper-proof voting.
//     - Add zero-knowledge proofs for anonymous voting.

// - Exists two types of transactions
//   1. Create new voting
//     - title
//     - voteId
//   2. Cast vote
//     - voteId

func main() {
}
