# h132: Experimental Secret Information Management System

[![Build Status][gh-actions-badge]][gh-actions]
[![go report][go-report-badge]][go-report]

**h132** is an experimental secret information management system written in Go.
It leverages affordable, tamper-resistant hardware modules such as **TPM 2.0** and **FIDO2 Security Keys** to provide a cost-effective alternative to traditional Hardware Security Modules (HSMs).

The name "h132" is derived from the Japanese word for "secret" (**秘密 - himitsu**) 😉

## ⚠⚠⚠️ Important Notice ⚠⚠⚠

This project is highly experimental. Before considering its use, it is *strongly recommended* to explore established tools like [GnuPG](https://gnupg.org/) and [age](https://age-encryption.org/). While these tools do not support the proposed multi-factor encryption scheme, they are arguably much safer than relying on experimental software.

## ✨ Key Features

- **Cost-Effective Solution:** Operates on a Raspberry Pi with a total cost under 20,000 JPY.
- **Tamper-Resistant Security:** Employs TPM 2.0 and FIDO2 Security Keys to ensure robust protection of encryption keys.
- **Secure Key Management:**
  - Generates encryption keys directly on the TPM 2.0 module.
  - Access keys are **not stored** on any devices; instead, they are dynamically generated using key material derived from the [HMAC-secret](https://fidoalliance.org/specs/fido-v2.1-ps-20210615/fido-client-to-authenticator-protocol-v2.1-ps-errata-20220621.html#sctn-hmac-secret-extension) via FIDO2 Security Keys.
- **Multi-Factor Secret Recovery:** Requires the following components to recover secrets:
  1. An **encrypted repository**.
  2. A hardware module connected to **TPM 2.0**.
  3. A **FIDO2 Security Key**.
  4. The **PIN** for the FIDO2 Security Key.
- **Emergency Recovery Mechanism:** Provides a "recovery phrase" to derive an emergency decryption key in case of hardware failure.

## 🛠️ How It Works

### At a High Level

1. **Access Key Generation:**
   - Access keys are generated and securely stored within the TPM 2.0 module.
2. **Access Secret Derivation:**
   - A salt value is sent to the FIDO2 Security Key, which generates an HMAC-secret.
   - This HMAC-secret is processed through a key derivation function to produce the access secret.
3. **Symmetric Encryption Key Retrieval:**
   - The access secret is used to decrypt the symmetric encryption key, which encrypts the file content.

### More Detail

*To be written.*

## 🖥️ Hardware Requirements

- **Raspberry Pi 4 (or any hardware that runs Linux)**
    - Due to the nature of this system, it is highly recommended to set up a dedicated Linux environment.
- **TPM 2.0 Module**
    - The OPTIGA™ TPM 2.0 Explorer clone module [GeeekPi Raspberry Pi TPM2.0モジュール TPM9670](https://amzn.asia/d/0k4sycS) is available for under 3,400 JPY.
    - You can use a TPM emulator like [swtpm](https://github.com/stefanberger/swtpm), but doing so would defeat the purpose of the system.
- **FIDO2 Security Key that supports the `hmac-secret` extension**
    - The [Yubico Security Key](https://www.yubico.com/products/security-key/) is available for under 7,000 JPY.
    - Be careful when purchasing a Security Key. Not all FIDO2 keys support the `hmac-secret` extension.

## 🚀 Getting Started

*To be written.*

## 🧑‍💻 Setting Up the Development Environment

For your convenience, devcontainer configuration files are provided. It is recommended to use them to set up your development environment.

### Using Dev Containers

To open the Project in a Dev Container:
- Open the `h132` project folder in VS Code.
- You should see a prompt asking to reopen the folder in a dev container. Click "Reopen in Container".
   - This will require "Remote Development" extension pack.
- VS Code will build and start the dev container defined by the provided configuration files.

### Using the TPM Emulator for Development

For development purposes, you can use the `swtpm` TPM emulator:

1. **Start the TPM Emulator:**
   - Run the script `hack/start_swtpm.sh` to start the TPM emulator.
2. **Set the TPM Path Environment Variable:**
   - Set the environment variable `H132_TPM_PATH` to point to the emulator socket:
     ```bash
     export H132_TPM_PATH=/workspaces/h132/var/swtpm_dir/server.sock
     ```
   - This tells the `h132` command to use the emulated TPM instead of a real device.

*Note:* While using a TPM emulator is acceptable for development, it defeats the purpose of the system for production use. Always use a real TPM 2.0 module in production environments.

## 📄 License
[Apache License 2.0](./LICENSE)

<!-- Markdown link & img dfn's -->
[go-report-badge]: https://goreportcard.com/badge/github.com/IPA-CyberLab/h132
[go-report]: https://goreportcard.com/report/github.com/IPA-CyberLab/h132
[gh-actions-badge]: https://github.com/IPA-CyberLab/h132/workflows/go/badge.svg
[gh-actions]: https://github.com/IPA-CyberLab/h132/actions