# KT Quick Start Guide

This guide will help you quickly set up KT (Kubernetes Troubleshooter) and start using it to diagnose issues in your Kubernetes clusters.

## 1. Installation

### 1.1 Prerequisites

Before installing KT, ensure you have the following prerequisites installed:

- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- kind (Kubernetes in Docker): [Install kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

### 1.2 Platform-specific Instructions

Choose your operating system:

#### Windows

1. Download the latest KT release for Windows from [GitHub Releases](https://github.com/your-kt-repo/releases).
2. Extract the ZIP file to a directory of your choice.
3. Add the directory to your system PATH.

#### macOS

1. Install KT using Homebrew:
   ```
   brew install kt
   ```

#### Linux

1. Download the latest KT release for Linux from [GitHub Releases](https://github.com/your-kt-repo/releases).
2. Extract the tarball:
   ```
   tar -xzvf kt-linux-amd64.tar.gz
   ```
3. Move the binary to a directory in your PATH:
   ```
   sudo mv kt /usr/local/bin/
   ```

### 1.3 Verify Installation

To verify that KT is installed correctly, run:

## 2. Setting Up a Test Environment

### 2.1 Create a Kubernetes Cluster with kind

1. Create a new cluster:
   ```
   kind create cluster --name kt-test
   ```

2. Verify the cluster is running:
   ```
   kubectl cluster-info --context kind-kt-test
   ```

3. Install [cloudtty](https://cloudtty.github.io/cloudtty/)

## 3. Using KT to Troubleshoot Kubernetes

### 3.1 Common Issues

#### ErrorPullImage

1. Deploy a pod with an invalid image:
   ```
   kubectl run test-pod --image=nginx2:latest
   ```

2. Use KT to diagnose the issue:
   ```
   kt diagnose pod test-pod
   ```

3. KT will identify the ErrorPullImage issue and provide suggestions for resolution.

### 3.2 More Complex Cases

For more complex troubleshooting scenarios, refer to the [KT Documentation](https://kt-docs.example.com/advanced-troubleshooting).

## 4. Additional Features

### 4.1 KTConsole

KTConsole provides a web-based interface for KT. To use it:

1. Start KTConsole:
   ```
   kt console
   ```

2. Open your browser and navigate to `http://localhost:8080`.

3. To find nodes for NodePort services, add the label `node-role.kubernetes.io/control-plane=` to the desired nodes:
   ```
   kubectl label nodes <node-name> node-role.kubernetes.io/control-plane=
   ```

## 5. Cloud Provider Specific Setup

### 5.1 Amazon EKS

If you're using Amazon EKS:

1. Install the AWS CLI: [AWS CLI Installation Guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)

2. Configure AWS CLI with your credentials:
   ```
   aws configure
   ```

3. Install eksctl: [eksctl Installation Guide](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html)

4. Create an EKS cluster:
   ```
   eksctl create cluster --name kt-eks-cluster --region us-west-2
   ```

5. Configure kubectl to use the EKS cluster:
   ```
   aws eks get-token --cluster-name kt-eks-cluster | kubectl apply -f -
   ```

Now you can use KT with your EKS cluster just like with any other Kubernetes cluster.

## 6. Next Steps

- Explore the [KT Documentation](https://kt-docs.example.com) for advanced usage and features.
- Join the [KT Community](https://kt-community.example.com) for support and discussions.
- Contribute to KT on [GitHub](https://github.com/your-kt-repo).

Happy troubleshooting with KT!
KTConsole:
- add label `node-role.kubernetes.io/control-plane` to find node for nodePort svc
EKS
- install aws

