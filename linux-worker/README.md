# linux-worker
Linux worker for Rabbit CI.

It provides resque compatible goworkers for config extraction and builds.

# Contributing setup

(This assumes you already have Rabbit CI setup and running.)

- Make sure you have the [worker-utils](https://github.com/rabbit-ci/worker-utils)
binary in your PATH.
- Make sure go and [gopm](https://github.com/gpmgo/gopm) are installed.
- Make sure redis is installed and running.
- Clone this repo.
- Run `gopm get` and `gopm build` in the folder you cloned to.
- Start it with: `./output --queues=workers --interval=1 --concurrency=2 --use-number=true`
