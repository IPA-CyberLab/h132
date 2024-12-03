# h132: Experimental Secret Information Management System
<p align='center'>
<img src="https://github.com/IPA-CyberLab/h132/blob/master/website/logo.png?raw=true" />
</p>

**h132** is an experimental secret information management system written in Go.
It leverages affordable, tamper-resistant hardware modules such as **TPM 2.0** and **FIDO2 Security Keys** to provide a cost-effective alternative to traditional Hardware Security Modules (HSMs).

The name "h132" is derived from the Japanese word for "secret" (**秘密 - himitsu**) 😉

[![Build Status][gh-actions-badge]][gh-actions]
[![go report][go-report-badge]][go-report]

## ⚠️🚧⚠️ Important Notice ⚠️️🚧⚠️ ️

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

## 📘 Usage Guide for h132

This document provides a step-by-step guide on how to use **h132**, an experimental secret information management system. Follow the instructions below to securely manage your secret files.

### 📁 Setting Up a Secret Repository

First, set up a Source Control Management (SCM) repository (like Git) to store your secret files. Although all files managed by the SCM are stored encrypted, the repository activity is not encrypted. Therefore, it is recommended to set up a **private repository**.

### 📝 Generating a Configuration File

Next, generate a configuration file for the repository. This file stores the `h132` configuration that applies to all file operations within the repository. The repository-wide configuration is called a **letter writing set**.

Run the following commands to create a new letter writing set:

```sh
$ export LWS_DIR=/path/to/your/repo
$ h132 lws create --name your_lws_name
info    Letter writing set (name=your_lws_name) successfully created!
```

- Replace `/path/to/your/repo` with the path to your secret repository where you want to store your letter writing set.
- Replace `your_lws_name` with a name of your choice for the letter writing set.

**Note:** Setting the `LWS_DIR` environment variable defines the location where `h132` will store its configuration files for the letter writing set.

### 🔑 Adding Access Keys

#### Adding a TPM 2.0 Managed Access Key

To add an access key managed by TPM 2.0 and protected by a FIDO2 Security Key, use the following command:

```sh
$ h132 keys add --type webauthn_wrapped_tpm --name your_key_name --tpmKeyHandle 0x81008015
info    Using TPM device: /dev/tpmrm0
Please navigate to the following URL in your browser: https://ipa-cyberlab.github.io/h132/webauthn_bridge/#...

info    Successfully registered WebAuthn credential. Now acquiring PRF secret.
Please navigate to the following URL in your browser: https://ipa-cyberlab.github.io/h132/webauthn_bridge/#...

info    Successfully acquired PRF. Now provisioning the key on the TPM.
info    Successfully generated key: [key details] {Name: your_key_name, Type: WebauthnWrappedTpm, TpmKeyHandle: 0x81008015, ReflectorUrl: https://ipa-cyberlab.github.io/h132/webauthn_bridge/}
info    Updated letter writing set: key "your_key_name" added.
info    Successfully added the key to the letter writing set.
```

- Replace `your_key_name` with a name for your access key.
- The `--tpmKeyHandle` option specifies the TPM handle where the key will be stored. We recommend using a handle in the range `0x81008000` to `0x81008FFF`.
    - Please consult Section 2.3.1 of the ["Registry of Reserved TPM 2.0 Handles and Localities"](https://trustedcomputinggroup.org/wp-content/uploads/RegistryOfReservedTPM2HandlesAndLocalities_v1p1_pub.pdf) for precise ranges.
    - `h132` will safely exit if there's a pre-existing key at the specified handle, so feel free to experiment.
- By [default](https://github.com/IPA-CyberLab/h132/blob/master/tpm2/device.go#L13), `h132` uses the TPM device at `/dev/tpmrm0`. If your TPM device is at a different path, specify it using the `H132_TPM_PATH` environment variable:
    ```sh
    $ export H132_TPM_PATH=/dev/tpm0
    ```

#### Adding an Emergency Access Key

To add an emergency access key (a backup method to access your secrets), run:

```sh
$ h132 keys add --type emergency --hint your_hint --name emergency
Mnemonic accepted:
[Your mnemonic phrase will be displayed here]

info    Successfully generated key: [key details] {Name: emergency, Type: EmergencyBackupKey, Hint: your_hint}
info    Updated letter writing set: key "emergency" added.
info    Successfully added the key to the letter writing set.
```

- Replace `your_hint` with a hint for where you will store your mnemonic phrase (e.g., "safety deposit box").
- The mnemonic phrase displayed is your recovery phrase. **Store it securely**.

### 📜 Listing Registered Keys

To list the keys registered in your letter writing set, use:

```sh
$ h132 keys ls
info    Found 2 keys in the letter writing set "your_lws_name"
[0] [key details] {Name: emergency, Type: EmergencyBackupKey, Hint: your_hint}
[1] [key details] {Name: your_key_name, Type: WebauthnWrappedTpm, TpmKeyHandle: 0x81008015, ReflectorUrl: https://ipa-cyberlab.github.io/h132/webauthn_bridge/}
```

### 📥 Importing a Sensitive File

To import (encrypt) a sensitive file into your repository, run:

```sh
$ h132 envelope seal --key your_key_name /path/to/your_secret_file
info    Using TPM device: /dev/tpmrm0
Please navigate to the following URL in your browser: https://ipa-cyberlab.github.io/h132/webauthn_bridge/#...

info    Successfully sealed h132 envelope "/path/to/your_secret_file.h132"
info    Produced an envelope file "/path/to/your_secret_file.h132"
```

- Replace `your_key_name` with the name of the access key you added earlier.
- Replace `/path/to/your_secret_file` with the path to the file you want to encrypt.
- The encrypted file will be saved with a `.h132` extension in `LWS_DIR` directory.

### 📤 Decrypting a File

To decrypt the file when you need to access its contents, use:

```sh
$ h132 envelope unseal --key your_key_name /path/to/your_secret_file.h132
info    Using TPM device: /dev/tpmrm0
Please navigate to the following URL in your browser: https://ipa-cyberlab.github.io/h132/webauthn_bridge/#...

info    Successfully unsealed h132 envelope "/path/to/your_secret_file.h132"
info    Produced a plaintext file "/path/to/your_secret_file"
```

- Replace `/path/to/your_secret_file.h132` with the path to the encrypted file.
- The decrypted file will be restored without the `.h132` extension in `LWS_DIR` directory.

---

**Note:** Throughout the process, you will be prompted to navigate to a URL in your browser. This is part of the WebAuthn authentication process using your FIDO2 Security Key. Follow the instructions in your browser to complete the authentication.

Remember to keep your TPM device, FIDO2 Security Key, and mnemonic phrase secure at all times to maintain the integrity of your secret management system.

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

### View envelope structure

You can use `h132 envelope dump` command to view `.h132` envelope file details:

```sh
$ h132 envelope dump var/lws2/password.txt.h132
Envelope:
  Signature: 304402200364ab382e3302b004f5ae35fd0237989a7a720c5e589dafcae312cb722b46790220612be44ae08e8c4bb85473926809692aa3a59e737320f1d2a5555bc6de367a6c (valid)
  Letter:
    len(Ciphertext): 24
    Nonce: 41f937a06bbf292bb894ea58
    RecipientKey[0]:
      RecipientPublicKey: xqcH32ShbyI+Ekr4Bm1K89LiIbUgHXR78cSvQbA5IWs= (emergency)
      SenderEphemeralPublicKey: bZn1pYXn/7aIHD5vnKbpZ6vvaGGI41FUydrP5NTNNKA=
      SenderEphemeralPublicKeySign: 304402207dc407c0198cfe34eb10f8faf72363859ecee47ea90cdc8b769872d209dd217e022050ffc07faca8ed53e76c182610afef5999408236cb144b98255ea7ef4baa4d5b (valid)
      EncryptedSymmetricKey: 7c1fb9abb83a71a18b4f9d658ee2def2d1ecbc1c80c5854e9151b1714f54b549027f327ebd02d3675d23e2e1d7391f98
      Nonce: 8d6d1a954100ac6d268fb4fc
    RecipientKey[1]:
      RecipientPublicKey: W7pd/FqIr8AK3ubizKmjqArN/f4yGLkmjXMxrlU/Ufg= (wallet_sk2)
      SenderEphemeralPublicKey: O94s7hXBJVuhUFL2dqA9HfSc9GCx+F8JShaUJvZbw/o=
      SenderEphemeralPublicKeySign: 3044022042773e467936fa08582b7a4ad15da8d87a849cbb79ae4af269bb5208591eb8a602202656bd711be13bc053f536e892d07b48823c31c183a90cf34b905641bb174e4a (valid)
      EncryptedSymmetricKey: 2092b02d610e6acc646c3e7010bf9ec3e1793f3c92d3b5c788a30c9fb841ad54e16a8a817fefb1138fb65c318c7603d9
      Nonce: 788f61f82385c9e74109a482
    SenderPublicKey: W7pd/FqIr8AK3ubizKmjqArN/f4yGLkmjXMxrlU/Ufg= (wallet_sk2)
```

## 📄 License
[Apache License 2.0](./LICENSE)

<!-- Markdown link & img dfn's -->
[go-report-badge]: https://goreportcard.com/badge/github.com/IPA-CyberLab/h132
[go-report]: https://goreportcard.com/report/github.com/IPA-CyberLab/h132
[gh-actions-badge]: https://github.com/IPA-CyberLab/h132/workflows/go/badge.svg
[gh-actions]: https://github.com/IPA-CyberLab/h132/actions