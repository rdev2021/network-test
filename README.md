# Introduction
 Network Test UI provides a user-friendly, web-based interface for essential network diagnostics, allowing users to Check connectivity, ports, HTTP status, domains, and databases without needing CLI expertise.

 This was a personal project aimed at assessing connectivity across subnetworks once firewall access was authorized and implemented. The container image makes it easy to deploy it on managed services such as EKS, AKS, App Services & ECS.

# Pre requisites
 You can run it directly on your local machine either using Docker/Podman or by executing Go commands.

# How to Run
 
 ## 1. Go
  ```bash
  go run main.go
  ```
 ## 2. Docker
  ```bash
  docker run -it --rm  -p 9091:9091 ghcr.io/rdev2021/network-test:latest
  ```
  http://localhost:9091/home/ to launch the home page of the application.

 ## 3. Kubernetes (Selfmanaged, GKE, AKS or EKS)

  The sample below shows NodePort services however use proper Ingress in an actual K8s environment.

  ```yaml
  ---
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: network-test
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: network-test
    template:
      metadata:
        labels:
          app: network-test
      spec:
        containers:
        - name: network-test
          image: ghcr.io/rdev2021/network-test:latest
          ports:
          - containerPort: 80

  ---
  apiVersion: v1
  kind: Service
  metadata:
    name: network-test-service
  spec:
    type: NodePort
    selector:
      app: network-test
    ports:
    - port: 80
      targetPort: 80
      nodePort: 30007 # Change to desired NodePort in the range 30000-32767
  ```

# Screenshots
 ### Home Page
 ![Screenshot](assets/nt-1.png)
 ### Port Check
 ![Screenshot](assets/nt-2.png)
 ![Screenshot](assets/nt-3.png)
 ### Http checks
 ![Screenshot](assets/nt-4.png)
 ![Screenshot](assets/nt-5.png)
 ### Domain checks
 ![Screenshot](assets/nt-6.png)
  ### Database checks
 ![Screenshot](assets/nt-7.png)

