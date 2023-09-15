<p align="center">
  <img src="/assets/punq_logo.png" alt="punq logo" width="400"/>
</p>

<p align="center">
    <a href="https://github.com/mogenius/punq/blob/main/LICENSE">
        <img alt="GitHub License" src="https://img.shields.io/github/license/mogenius/punq?logo=GitHub&style=flat-square">
    </a>
    <a href="https://github.com/mogenius/punq/releases/latest">
        <img alt="GitHub Latest Release" src="https://img.shields.io/github/v/release/mogenius/punq?logo=GitHub&style=flat-square">
    </a>
    <a href="https://github.com/mogenius/punq/releases">
      <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/mogenius/punq/total">
    </a>
    <a href="https://github.com/mogenius/punq">
      <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/mogenius/punq">
    </a>
</p>
<p align="center">
    <a href="https://github.com/mogenius/punq">
      <img alt="GitHub repo size" src="https://github.com/mogenius/punq/actions/workflows/main.yml/badge.svg">
    </a>
    <a href="https://github.com/mogenius/punq">
      <img alt="GitHub repo size" src="https://github.com/mogenius/punq/actions/workflows/develop.yml/badge.svg?branch=develop">
    </a>
</p>

# punq

punq provides a WebApp and CLI to easily manage multiple Kubernetes clusters. It comes with integrated team collaboration, logs, and workload editor for clusters across different infrastructures. The goal of this project is to make DevOps' lifes easier by improving Kubernetes operations especially in teams.

## How it works

punq is self-hosted on a Kubernetes cluster to run an instance for you and your team members. Each instance consists of the following services:

- The operator written in Golang
- An Angular application serving the user interface

With punq you can then manage multiple Kubernetes clusters by adding them from your local kubeconfig. The configurations are stored as secrets on your cluster and based on them punq displays all workloads and resources in the application. This way, every user of your punq instance can monitor and manage clusters without requiring access to the kubeconfig.

## Installation

The setup of punq is done via the command line interface. You can install it with the following commands.

### Mac/Linux

```
brew tap mogenius/punq
brew install punq
```

### Windows

```
Install: https://scoop.sh/
scoop bucket add mogenius https://github.com/mogenius/punq
scoop install punq
```

## Setting up punq

Once you installed the punq CLI here's how to set it up.

```
# List all available CLI features
punq help
# Install the punq operator in your current kubecontext
punq install
# start punq on the cluster
punq
# Set up the ingress with your domain to serve the punq web application
punq -i yourdomain.com
# Manage users and permissions
punq user
```

The admin credentials for your punq instance are prompted to your terminal when punq is started. Use them to log in to the punq web application and start adding clusters.

Have fun with punq! 🤘

## Development

To update the documentation please run (in project root):

```
swag init --parseDependency --parseInternal
```

## Contribution

punq is still still early stage and we're inviting you to contribute. Feel free to pick up open issues and create PRs.

---

Made with 💜 by the folks from [mogenius](https://mogenius.com)

#
