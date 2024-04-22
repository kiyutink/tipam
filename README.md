# Tipam
Tipam is an IP address manager for the terminal. It solves a problem of managing a hierarchical structure of networks, and validating new networks on creation.


## Main concepts
**Claim** - The most important. A claim is a record of 3 fields:
  - *CIDR* - The CIDR notation of a network
  - *Tags* - List of strings
  - *Final* - A boolean field

The entire hierarchy of the networks is laid out using *claims*. For example, a claim that looks like this:
```
CIDR: 10.0.0.0/24
Tags: staging
Final: false
```
Specifies that the network with CIDR 10.0.0.0/24 is claimed (reserved) for the staging environment. At the same time it's not a *final* claim, which means that you can create **subclaims** of it.

**Subclaim** - A claim that specifies a *smaller* network that fits within some other claim's network. For example, for the claim above, we can create a subclaim like this:
```
CIDR: 10.0.0.0/28
Tags: staging / website
Final: false
```
This claim is for a smaller network (10.0.0.0/28), and is *valid* with respect to its **superclaim** (10.0.0.0/24). A subclaim is *valid* if it's *tagged* with all the same tags as its direct superclaim, and at least one additional tag thereafter. E.g. `staging/website` or `staging/website/europe` or `staging/test` are all valid with respect to a superclaim tagged with `staging`. `production` or `test/staging` or `staging` are all invalid with respect to the same superclaim.

**Superclaim** - For a **subclaim**, each larger claim that it fits in is a **superclaim**

**Persistor** - The way to persist state. By default state is persisted in a local yaml. You can also use an `inmemory` persistor (for debugging) and `s3dynamo` to keep the state in an s3 bucket and lock access to it with Dynamo DB. Basically same thing as `backend` in terraform land.

## Usage
Tipam can be used in one of 3 ways:
- As a regular CLI (e.g. running commands like `tipam claim --cidr ... --tag ...`)
- As a Terminal UI application (e.g. running `tipam`)
- As a go package (installable with `go get` and `import "tipam/tipam"`)
