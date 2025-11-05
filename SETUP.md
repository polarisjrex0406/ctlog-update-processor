# SETUP

## Environment
- OS: Ubuntu 24.04.3 LTS

## PostgreSQL Installation

### 1. Add PostgreSQL Repository
```sh
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget -qO- https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
```

### 2. Update Package List
```sh
sudo apt update
```

### 3. Install PostgreSQL 16
```sh
sudo apt install postgresql-16
```

### 4. Verify Installation
```sh
sudo systemctl status postgresql
psql --version  # Should show e.g., 16.8
```

### 5. Post-Install Configuration

1. **Set Password** for postgres user

```sh
sudo -u postgres psql -c "ALTER USER postgres PASSWORD '123456aA@';"
```

2. **Enable Remote Access**

```sh
sudo nano /etc/postgresql/16/main/postgresql.conf
```

Uncomment/edit:

```ini
listen_addresses = '*'
```

3. **Upate Access Rules:**

```sh
sudo nano /etc/postgresql/16/main/pg_hba.conf
```

Add line:

```ini
host    all             all             0.0.0.0/0               md5
```

### 6. Restart Service

```sh
sudo systemctl restart postgresql
```

## Add SSH Key to BitBucket

### 1. Check Existing SSH Keys

Run:

```sh
ls -al ~/.ssh
```

Look for files like `id_ed25519.pub` or `id_rsa.pub`. If missing, generate a new key:

### 2. Generate SSH Key (if needed)

```sh
ssh-keygen -t ed25519 -C "your_email@example.com"
```

Press Enter to accept default paths. Add a passphrase if desired.

### 3. Add Public Key to Bitbucket

- Display your public key:

```sh
cat ~/.ssh/id_ed25519.pub
```

- Copy the output (starts with `ssh-ed25519 ...`).

- Go to **Bitbucket** → **Settings** → **Personal Bitbucket settings** → **SSH keys** and paste it.

## Build and Install PostgreSQL Extension Libraries

### 1. Install Go 1.25.1

1. Download the Go Binary

```sh
wget https://golang.org/dl/go1.25.1.linux-amd64.tar.gz
```

2. Extract the Archive

```sh
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
```

3. Set Environment Variables

Add these lines to your shell profile (e.g., `~/.profile` or `~/.bashrc`):

```sh
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
```

Then reload the profile:

```sh
source ~/.profile  # or source ~/.bashrc
```

4. Verify Installation

```sh
go version
```

Expected output: `go version go1.25.1 linux/amd64`

### 2. Install PostgreSQL Server Development Files

```sh
sudo apt-get install postgresql-server-dev-16
```

### 3. Install make, gcc

```sh
sudo apt update
sudo apt install make gcc
```

### 4. libzlintpq

1. Download libzlintpq

```sh
git clone git@bitbucket.org:xoduxcrt/libzlintpq.git
```

2. Download PL/Go

```sh
cd libzlintpq
git clone https://gitlab.com/microo8/plgo.git
```

3. Build PL/Go

Build plgo

```sh
cd ..
go build -o $GOPATH/bin/plgo
```

Download plgo repository

```sh
mkdir $GOPATH/src/gitlab.com/microo8
cd $GOPATH/src/gitlab.com/microo8
git clone https://gitlab.com/microo8/plgo.git
cd plgo
```

Modify pl.go

```diff
Datum cstring_to_datum(char *val) {
-    return CStringGetDatum(cstring_to_text(val));
+    return PointerGetDatum(cstring_to_text(val));
}

char* datum_to_cstring(Datum val) {
-    return DatumGetCString(text_to_cstring((struct varlena *)val));
+    return text_to_cstring((text *) DatumGetPointer(val));
}

bytea* datum_to_byteap(Datum val) {
-    return DatumGetByteaPP((struct varlena *)val);
+    return DatumGetByteaP(val);
}

func (db *DB) Prepare(query string, types []string) (*Stmt, error) {
	var typeIds []C.Oid
	var typeIdsP *C.Oid
	if len(types) > 0 {
		typeIds = make([]C.Oid, len(types))
		var typmod C.int32
		for i, t := range types {
			ct := C.CString(t)
			defer C.free(unsafe.Pointer(ct))
-			C.parseTypeString(ct, &typeIds[i], &typmod, (C._Bool)(false))
+			C.parseTypeString(ct, &typeIds[i], &typmod, (*C.Node)(nil))
    ...
}
```

4. Build & Install libzlintpq

```sh
GOPATH=/root/go go get -u gitlab.com/microo8/plgo/...
GOPATH=/root/go go get -u github.com/zmap/zlint
CGO_ENABLED=1 make
cd build
make install with_llvm=no
```

### 5. libx509pq

1. Download libx509pq

```sh
git clone git@bitbucket.org:xoduxcrt/libx509pq.git
```
2. Build & Install libx509pq

```sh
make
make install
```

### 6. libocsppq

1. Download libocsppq

```sh
git clone git@bitbucket.org:xoduxcrt/libocsppq.git
```

2. Build & Install libocsppq

```sh
make
cd build
make install with_llvm=no
```

## Create Database Schema in PostgreSQL

### 1. Download certwatch_db into non-root directory

```sh
mkdir /etc/crtsh
cd /etc/crtsh
git clone git@bitbucket.org:xoduxcrt/certwatch_db.git
```

### 2. Execute SQL

```sh
cd certwatch_db
sudo -u postgres psql
\i sql/create_schema.sql
```

## Build CT Log Scraper

### 1. Install Rust

1. Update System Packages

```sh
sudo apt update && sudo apt upgrade -y
```

2. Install Required Dependencies

```sh
sudo apt install curl build-essential
```

3. Download and Run Rust Installer

```sh
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

- Press `1` to proceed with default installation.

4. Load Environment Variables

```sh
source $HOME/.cargo/env
```

(Alternatively, restart your terminal)

5. Verify Installation

```sh
rustc --version  # Should show version like "rustc 1.xx.x"
cargo --version  # Rust's package manager
```

### 2. Build scrape-ct-log

1. Download mpalmer/scrape-ct-log

```sh
git clone https://github.com/mpalmer/scrape-ct-log.git
```

2. Build mpalmer/scrape-ct-log

```sh
cd scrape-ct-log
cargo build --release
```

You'll end up with a binary at `target/release/scrape-ct-log`.

Make a directory named `tmp` in `target/release`.

```sh
mkdir target/release/tmp
```

## Download and Run Update Processor

### 1. Download ctlog-update-processor

```sh
git clone git@bitbucket.org:xoduxcrt/ctlog-update-processor.git
```

### 2. Run ctlog-update-processor

1. Initialize CT Log Config

```diff
main.go
...
+   utils.ScrapeCTLogList()
+   certwatch.InitConfig()
...
```

2. Run ctlog-update-processor

```diff
main.go
...
+   // utils.ScrapeCTLogList()
+   // certwatch.InitConfig()
...
```