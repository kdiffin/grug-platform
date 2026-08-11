# deprecated 
this project is for ease of deployments in K8s

it gives you a nice CLI interface for working with k8s native deployments.

You will be able to easily manage stuff like:

Databases/queues/cache with mandatory backups
Built in FOSS observability with grafana / whatever I want to choose 
Built in argoCD / gitops integration
and so on..

The configs and what it installs in terms of operators will be transparent.

We're not trying to make a new platform like AWS, we're just taking the existing building blocks that allow on-prem cloud native deployments to work and making it easy to setup.

Theres only so many things apps and infra needs, I think those are mainly:

Version control (git does this already, no gap here)

Backups restorability, and rollbacks (ik how id do this for cloudnativepg, i have to look at the others and figure out whats the best practice for those though and really get that dialed in.) 

Also I should look at how to do this on the application level. stateful apps should be thought through (task) and how devs request disk space should be thought through. 

I believe we can add some sort of stanza for this.  


## interface 

Grug in the future can have these commands:

`grug spin db` -> gives a list of dbs to choose from, and their sizes, s, m, l or custom 
`grug spin cache` -> redis, memcached 
`grug spin message-queue` -> 
`grug deploy agent` -> deploys an agent which can be configured to be any agent you want (hermes, openclaw, etc.) with the permissions you want 

and these will automatically provide backups for these projects (so for example if its cloudnativepg, it will back it up to either a local s3 (if the uri is provided) or aws s3.) spinning up the db without backups should be allowed but should be heavily advised against in the CLI.

## personal thoughts

I think the program should also have primitives to spin up a k8s cluster easily, but like, on steroids though.

Just run grug init, and boom. Everything is setup.

In the future this can be done via integrations with AWS, GCP, Azure, etc.

For the initial scope now though im just going to get those commands working with sensible defaults. Kafka with streamzi, Postgres with cloudnativepg, easy argocd / grafna integration, automatic backups, etc.

And then we can integrate user applications to this easily as well. maybe link it to a git repo with a simple yaml which defines which resources it wants, and grug handles it itself and shows it in the cli when running 

`grug list (resource type)`

`note-to-self` I should be wary of not recreating k8s.
