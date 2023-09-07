🛠 Simple Proof of Work Blockchain with EVM

This project provides a minimalistic simulation of a blockchain system that includes a rudimentary implementation of the Ethereum Virtual Machine (EVM). It's designed as an educational tool to help understand the core concepts of blockchain and the EVM.
✨ Features:

    🔗 Proof of Work (PoW) Consensus Algorithm
    📦 Basic Block Structure
        📝 Hash
        ⛓ Previous Block Hash
        🔄 Transactions
        🌲 State Root
        🎲 Nonce
    💸 Simple Transactions Model
    🖥 Rudimentary Ethereum Virtual Machine (EVM) Implementation
        🚀 PUSH
        ➕ ADD
        ➖ SUB
        💾 SSTORE
        🔄 SLOAD

🚀 Quick Start:

    🔍 Clone the Repository

    bash

git clone [https://github.com/rohallah12/POW_EVM_Chain_Simple]

📂 Navigate to Project Directory

bash

cd [https://github.com/rohallah12/POW_EVM_Chain_Simple]

▶️ Run the Program

bash

    go run main.go

💡 How It Works:

    📝 Creating a Transaction: Utilize the CreateTransferTx function to generate a transaction that transfers tokens from one address to another.

    🔗 EVM Operations: The basic EVM simulation can execute simple operations, giving a glimpse into how Ethereum processes instructions in its contracts.

    ⛏️ Mining a Block: The MineBlock function showcases a simplified proof of work mechanism, where the system calculates hashes until it finds one below a set difficulty.

    🛡️ Verifying Block's Proof of Work: The verify function can be used to ensure that a block's hash meets the required proof of work.

🤝 Contribution:

Feel free to fork this repository and extend it in any way you see fit. If you have any ideas or suggestions, please open an issue or pull request.
⚠️ Disclaimer:

This project is for educational purposes only and does not serve as a full-fledged blockchain or EVM implementation.