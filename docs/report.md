# Lab Report
Group number: 12 <br>
Group members: Matteo Carnelos, Fernando Labra Caso, Fernando Castell <br>
Code repository: https://github.com/matteocarnelos/kadlab

## Introduction

This document contains an explanation of the process of developing a _Kademlia Distributed Hash Table (DHT)_ using the 
Go programming language. The project has been developed for the **Mobile and Distributed Computing Systems** course's 
laboratory at LTU.

Whether we choose one tool instead of another or we decide to use some specific methodology, everything will be 
recorded.

## Frameworks and Tools

### Markdown
We use Markdown as the markup language to produce the Lab Report of this project. It is as elegant and easy to use as 
it is efficient in providing structured documents. Moreover, it is completely integrated within the GitHub environment. 

Another option for writing this document was LaTeX, which at a starting point we had considered, but due to a more 
complex syntax we decided to leave it aside. If in the future we will need some LaTeX features that Markdown is not able
to provide, we might consider switching.

In particular, we are using the GitHub-flavored Markdown and the Pandoc tool to convert the markdown file into PDF.

### Go Language
As it was recommended by the professor, we decided to use Golang for the development of this project, not only because 
of its appealing syntax and the first class object usage, but also because of its management of threads and channels, 
which comes in handy for distributed programming.

At the time of writing this document, we are using Go 1.17, but we might update to newer versions later in the project.
The project is structured as a Go module, with the following directory structure:
```
.
├── docs/
├── kademlia/
├── test/
├── go.mod
└── main.go
```
With the `main.go` being the entrypoint of the module, the `docs` directory containing the documentation, the `kademlia`
directory containing the Kademlia related files, the `test` directory containing the test scripts and the `go.mod` file
defining the module properties (name, dependencies,...).

To write and manage the project we decided to use an Integrated Development Environment (IDE). In particular, we decided
to use the GoLand IDE, developed by JetBrains.

### Git
We are using Git as the Version Control System (VCS) for the project as we have a common repository with the codebase, 
and we need a tool to manage the different contributions.

In particular, we decided to use the Gitflow Workflow to standardize the development process through VCS, using the 
following branching structure:
 * Single `main` branch containing tested and production-ready code;
 * Single `develop` branch containing tested but not production-ready code;
 * Multiple `feature/...` branches containing new features in the development/testing phase that will be merged into the
   `develop` branch;
 * Multiple `release/...` branches containing tested and production-ready code that will be merged into the `main`
   branch;
 * Multiple `hotifx/...` branches containing bug fixes for the `main` branch.

To manage the repository more easily we are using the Fork git client, which comes with all the integrations needed for 
the Gitflow Workflow and the GitHub platform.

### GitHub
We use the GitHub platform to host and manage the repository and all the related components. It is a common choice for 
the storage of code, and we all had used it in other university projects, so we considered it appropriate.

In particular, we are using the following features of the GitHub platform:
 * Issues, Pull Requests, Code Reviews,... to easily manage the development workflow;
 * GitHub Actions for Continuous Integration (CI) and Continuous Deployment (CD);
 * GitHub Container Registry for hosting custom Docker images.

As for the CI/CD, we set up two workflows:
 * Docker CD: to automatically build and publish the Docker image when new code is pushed on the `main` branch;
 * Pandoc CD: to automatically generate the PDF report when the markdown version is updated.

### Docker
As suggested, we decided to use Docker to containerize our application and replicate them in a common environment using 
Docker swarm mode. 

We built our container on top of the `golang` base image, that can be found in the Docker Hub Container Image Library.
To replicate the containers and connect them in a network we used the `docker stack deploy` command, after enabling the 
Docker swarm mode.

Although the Docker CLI is intuitive, we decided to use a tool to graphically manage container deployment and 
management, namely, Portainer. With Portainer, we are able to jointly set up and manage the network throughout the team. 

### LUDD Distributed Systems and Technologies
For the deployment of the container, we needed a virtual machine that could hold up to as many instances as required for
the distributed system to work in the best condition.

As LTU students, LUDD (Luleå Academic Computer Society) offers services in the computing field, and their subsection
DUST (Distributed Systems and Technologies) virtual machines with different configurations, which we certainly chose for
deploying our containers and code.

We set up a Virtual Machine (VM) with the following specifications:
 * Operating System: Ubuntu 20.04 LTS
 * vCPUs: 8
 * RAM: 16384 MB
 * Disk: 64 GB

Finally, we enabled the Public IP support in order to remotely access the VM and login to the Portainer dashboard. 

## System Architecture
_[To be defined]_

## Limitations
_[To be defined]_

## Conclusion
_[To be defined]_

---

## Sprints

### Sprint 0

### Sprint 1

### Sprint 2
