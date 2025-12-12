# EC2 Self-Hosted GitHub Actions Runner Setup

This guide sets up a GitHub Actions self-hosted runner on an EC2 instance with Docker running as a **non-root user** (security best practice).

---

## Prerequisites

- Ubuntu EC2 instance (Amazon Linux works with minor adjustments)
- SSH access to the instance
- GitHub repository access to generate runner token

---

## Step 1: Install Docker

```bash
sudo apt-get update
sudo apt-get install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker
```

Verify installation:

```bash
docker --version
```

---

## Step 2: Create a Dedicated Runner User

```bash
# Create user with home directory and bash shell
sudo useradd -m -s /bin/bash github-runner

# Add user to docker group (allows Docker commands without sudo)
sudo usermod -aG docker github-runner
```

---

## Step 3: Install GitHub Actions Runner

### 3.1 Switch to the runner user

```bash
sudo su - github-runner
```

### 3.2 Download and extract the runner

```bash
# Create a folder
mkdir actions-runner && cd actions-runner

# Download the latest runner package (check GitHub for latest version)
curl -o actions-runner-linux-x64-2.329.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.329.0/actions-runner-linux-x64-2.329.0.tar.gz

# Optional: Validate the hash
echo "194f1e1e4bd02f80b7e9633fc546084d8d4e19f3928a324d512ea53430102e1d  actions-runner-linux-x64-2.329.0.tar.gz" | shasum -a 256 -c

# Extract the installer
tar xzf ./actions-runner-linux-x64-2.329.0.tar.gz
```

### 3.3 Configure the runner

Get your token from: **GitHub Repo → Settings → Actions → Runners → New self-hosted runner**

```bash
./config.sh --url https://github.com/YOUR_USERNAME/YOUR_REPO --token YOUR_TOKEN
```

**Configuration prompts:**

| Prompt | Response |
|--------|----------|
| Runner group | Press Enter (Default) |
| Runner name | Press Enter or custom name (e.g., `ec2-runner`) |
| Labels | Press Enter (default: `self-hosted,Linux,X64`) |
| Work folder | Press Enter (default: `_work`) |

### 3.4 Exit back to sudo user

```bash
exit
```

---

## Step 4: Install as a Systemd Service

This ensures the runner starts automatically on reboot and keeps running after SSH disconnect.

**One-liner (install and start):**

```bash
sudo bash -c 'cd /home/github-runner/actions-runner && ./svc.sh install github-runner && ./svc.sh start'
```

**Or run separately:**

```bash
# Install the service
sudo /home/github-runner/actions-runner/svc.sh install github-runner

# Start the service
sudo /home/github-runner/actions-runner/svc.sh start

# Verify it's running
sudo /home/github-runner/actions-runner/svc.sh status
```

---

## Service Management Commands

```bash
# Check status
sudo /home/github-runner/actions-runner/svc.sh status

# Stop the runner
sudo /home/github-runner/actions-runner/svc.sh stop

# Start the runner
sudo /home/github-runner/actions-runner/svc.sh start

# Uninstall the service
sudo /home/github-runner/actions-runner/svc.sh uninstall
```

---

## Verify Docker Works for Runner User

```bash
sudo -u github-runner docker run hello-world
```

---

## Troubleshooting

### Runner not picking up jobs

```bash
# Check service status
sudo /home/github-runner/actions-runner/svc.sh status

# Check logs
sudo journalctl -u actions.runner.* -f
```

### Docker permission denied

```bash
# Ensure user is in docker group
groups github-runner

# If docker group missing, add it and restart service
sudo usermod -aG docker github-runner
sudo /home/github-runner/actions-runner/svc.sh restart
```

### Token expired

Tokens are single-use and expire quickly. Generate a new one from GitHub:
**Settings → Actions → Runners → New self-hosted runner**

---

## Security Notes

| Approach | Security |
|----------|----------|
| `chmod 666 /var/run/docker.sock` | ❌ Insecure - any user can control Docker |
| `sudo docker ...` | ❌ Requires root privileges |
| **User in `docker` group** | ✅ Secure - only authorized users |

---

## Quick Reference (Copy-Paste)

```bash
# Full setup in one go (adjust token and repo URL)
sudo apt-get update && sudo apt-get install -y docker.io
sudo systemctl start docker && sudo systemctl enable docker
sudo useradd -m -s /bin/bash github-runner
sudo usermod -aG docker github-runner
sudo su - github-runner -c '
  mkdir -p ~/actions-runner && cd ~/actions-runner
  curl -o actions-runner-linux-x64-2.329.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.329.0/actions-runner-linux-x64-2.329.0.tar.gz
  tar xzf ./actions-runner-linux-x64-2.329.0.tar.gz
'
# Then run config manually:
# sudo su - github-runner
# cd actions-runner && ./config.sh --url https://github.com/OWNER/REPO --token TOKEN
# exit
# sudo /home/github-runner/actions-runner/svc.sh install github-runner
# sudo /home/github-runner/actions-runner/svc.sh start
```

