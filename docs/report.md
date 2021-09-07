# Lab Report
Group number: 12 <br>
Group members: Matteo Carnelos, Fernando Labra Caso, Fernando Castell Miñón <br>
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
 * Single `main` branch containing tested and production-ready code
 * Single `develop` branch containing tested but not production-ready code
 * Multiple `feature/...` branches containing new features in the development/testing phase that will be merged into the
   `develop` branch
 * Multiple `release/...` branches containing tested and production-ready code that will be merged into the `main`
   branch
 * Multiple `hotfix/...` branches containing bug fixes for the `main` branch

To manage the repository more easily we are using the Fork git client, which comes with all the integrations needed for 
the Gitflow Workflow and the GitHub platform.

### GitHub
We use the GitHub platform to host and manage the repository and all the related components. It is a common choice for 
the storage of code, and we all had used it in other university projects, so we considered it appropriate.

In particular, we are using the following features of the GitHub platform:
 * Issues, Pull Requests, Code Reviews,... to easily manage the development workflow
 * GitHub Actions for Continuous Integration (CI) and Continuous Deployment (CD)
 * GitHub Container Registry for hosting custom Docker images

As for the CI/CD, we set up two workflows:
 * Docker CD: to automatically build and publish the Docker image when new code is pushed on the `main` branch
 * Pandoc CD: to automatically generate the PDF report when the markdown version is updated

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

The organization of the sprints will be divided by content. Each member of the group will take charge of a different 
objective, as usual in the world of technology companies. 

After an initial meetup, we decided to overall split the work in the following way: network formation will be held by 
Fernando Castell, object distribution will be carried out by Fernando Labra and node management and application 
interface will be done by Matteo Carnelos. Unit testing will be developed individually as the codebase grows. Last but 
not least, the report will be written as the project is progressing, meaning that every member of the group will have to
add his individual part of the report while making any progress in the project.

We now focus on each sprint, giving more detailed information.

### Sprint 0
Our first priority when we were assigned this project was clear: a proper kademlia comprehension was mandatory. With
that in mind, our first approach to the project was purely theoretical: reading kademlia paper, looking for more
information on the internet and searching for some interesting videos on the internet were our first steps. After that,
we had to make the work distribution, which will be presented later on this document. Afterwards, we started working to
spin up a network of at least 50 containerized nodes. As presented on the paper, "the nodes do not have to carry any
Kademlia-related software at this point", making the containerization section simpler in a way.

#### Plan
 * Individual and collective study of the Kademlia algorithm principles, applying them with the help of simplified 
   examples
 * Sprint planning and job partitioning
 * Environment architecture planning, setup and testing
 * Report writing and reviewing
 * Demo code production

#### Backlog
 * Completion of the first five mandatory objectives:
   * M1 - Network formation: implementation of pinging, network joining and node lookup
   * M2 - Object distribution: add to the nodes the functionality to store (and to find) data objects
   * M3 - Command line interface: a CLI must be added to be able to execute the put, get and exit commands
   * M4 - Unit testing: add tests to check the proper functioning of the code
   * M5 - Containerization: deploy a network of at least 50 nodes in the same virtual machine with each node in a container
 
 * Completion of the qualifying objectives:
   * U1 - Object expiration: the TTL mechanism should be used to limit the lifetime of data in the network. It will also be
     necessary to decide the TTL, as if it is changeable or not
   * U2 - Object expiration delay: to avoid losing the objects more than one node should have the information; in this 
     objective, the main goal is to keep contact with those nodes, so they do not delete the information.
   * U3 - Forget CLI command: allow the original uploader of an object to stop refreshing it.
   * U4 - RESTful application interface: provide compatibility with web applications by implementing the RESTful API
     on each node and thus being able to transfer the files using HTTP methods.
   * U5 - Higher unit test coverage: add even more complete tests to check the proper functioning of the code
   * U6 - Concurrency and thread safety: make use of the golang utilities (channels or threads) for sending concurrent messages
          to the other nodes
 
 * Update the lab report

#### Reflections
The initial project is based on the development of a Distributed Data Store (DDS), but as the implementation of the Kademlia software 
is not required at this stage, the overall reflection is that this first nil sprint was not a challenge in terms of code but more 
in terms of writing this report and understanding the functioning of Kademlia.

### Sprint 1
#### Plan
* Completion of the first five mandatory objectives, the order applied in the lab project 
  explanation is suitable for the priority development of our code:
    1. Network formation: implementation of pinging, network joining and node lookup
    2. Object distribution: add to the nodes the functionality to store (and to find) data objects
    3. Command line interface: a CLI must be added to be able to execute the put, get and exit commands.
    4. Unit testing: add tests to check the proper functioning of the code.
    5. Containerization: automize the start and stop of the network using a script.
* Sprint 2 planning
* Update lab report

#### Backlog
#### Reflections 

### Sprint 2
#### Plan
#### Backlog
#### Reflections 
