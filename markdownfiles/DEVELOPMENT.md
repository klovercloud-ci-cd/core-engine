### Preferred OS

- ```Linux```
- ```Mac```

### Required tools

- [.git](https://docs.github.com/en/get-started/quickstart/set-up-git).
- [golang](https://go.dev/doc/install).
    - ``Note`` Golang version ```1.16``` or higher is recommended.
- [Docker](https://docs.docker.com/engine/install/ubuntu/)
- kubernetes.
    - ``Note`` version ```1.22.2``` or higher is recommended. In case of ```minikube```, you may
      try ```minikube start --kubernetes-version=v1.22.2```
- [tekton](https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.34.1/release.yaml)

### Preferred IDE

- [goland](https://www.jetbrains.com/go/)

### Prepare Environment

- Clone source code from you forked repository.
- Create ``.env`` file in project base directory
    - Find environment variables from ```.examle_env``` file
- Create ``.env.mongo.test`` file in project base directory
    - Find environment variables from ```.examle_env.mongo.test``` file