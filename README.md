# Orchid
Command line application and CI server for execute a defined set of jobs
locally and remotely through the same interface. The project focuses on simple
deployment and easy replication of everything  needed for deployment accross
computers. The application is written entirely in Go.

**Note: The project is under development. It is not finished, and the features described below can not be expected to work.**


# Usage
The application has the following interface:

```
- list jobs     // List all configured jobs
- list machines // List all configured machines
- list scripts  // List all configured scripts
- list logs     // List all stored logs
- run <job id>  // Run the job with the given id
- logs <log id> // Tail the log with the given id
```

It looks for a directory named `orchid` in which the configuration files reside
as described further below.


# Installation
Orchid requires docker to run. Clone this repository and add the `scripts`
directory to your path.


# Configuration
Everything in Orchid is defined in flat files, making an entire setup easily
versioned using Git and easily configured using any text editor.

Orchid project a very simple to setup and use. A setup consists of definition
of the following entities described in detail below:

- Machines
- Jobs
- Scripts
- Keys
- Server (optional)
- Logs

The configuration files are expected to reside in a directory named `ci` with
the following structure:

```
- jobs.json
- keys
--- <RSA private keys for SSH>
- logs
--- <Log files managed by Orchid>
- logs.json
- machines.json
- scripts
--- <Executable files>
```


## Machines
A machine is a remote server on which commands can be executed. This is useful
for deployment. A machine definition consists of the following attributes:

- **Id:** A unique machine identifier
- **Address:** The IP address / URL at which the machine resides
- **Port:** The SSH port used by the machine
- **User:** The username used for accessing the machine through SSH
- **PrivateKey:** The name of private key needed for accessing the machine
  through SSH (path to relative to the `keys` directory)

The configuration resides in the `machines.json` file. A sample config file is
given below:

```
[
  {
    "Id": "machine1",
    "Address": "192.168.1.50",
    "Port": "3333",
    "User": "myuser",
    "PrivateKey": "machine1.key"
  },
  {
    "Id": "machine2",
    "Address": "127.0.0.2",
    "Port": "1234",
    "User": "someuser",
    "PrivateKey": "machine2.key"
  }
]
```


## Jobs
A job is the unit of execution. A job definition consists of the following
attributes:

- **Id:** A unique job identifier
- **Pipeline:** A list of machine/script pairs to execute in the job:
    - **Machine:** Identifier of the machine on which to run the script or the
      value "local" indication that the script is executed locally
    - **Script** The name of the script / executable file to run (path relative
      to the `scripts` directory)

The configuration resides in the `jobs.json` file. A sample config file is
given below:

```
[
  {
    "Id": "job1",
    "Pipeline": [
      {
        "Machine": "local",
        "Script": "script1.sh"
      },
      {
        "Machine": "local",
        "Script": "script2.sh"
      }
    ]
  },
  {
    "Id": "job2",
    "Pipeline": [
      {
        "Machine": "machine2",
        "Script": "script1.sh"
      },
      {
        "Machine": "machine1",
        "Script": "script2.sh"
      }
    ]
  }
]
```


## Scripts
The concept of script covers the executable files located in the `scripts`
directory. These are the executables available in the job definitions.


## Keys
The concept of keys covers the RSA private keys located in the `keys`
directory. These are the keys used for accessing remote machines.


## Server (optional)
The server definition is needed if you wish to execute jobs on a running
instance of the Orchid CI server. Only a single server configuration is
allowed. A server definition consists of the following entities:

- **Address:** The IP address / URL at which the machine resides
- **Port:** The SSH port used by the machine
- **Secret:** Secret/token/key used for validating access to the server. This
  secret can also be used in Github Webhooks


## Logs
Logs are managed entirely by the Orchid application. Metadata about the logs is
stored in the `logs.json` file. The output of job executions are stored in
files in the `logs` directory.


# Installation
TODO


# Features to be implemented

- Getting the list of logs on remote machines
- Getting the log output on remote machines
- Log filtering when listing logs
- Log synchronization between the local and remote machine
