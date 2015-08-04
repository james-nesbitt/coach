# COACH

Coach is a docker-compose port to go, with an augmented yaml syntax.

Coach allows:
- defining some nodes as being of different types
- some node types can be used just for builds, or just for volumes
- some nodes can be used as utility commands

Coach handles instances slightly differently:
- some nodes can be considered "scalable"
- some nodes can have only a single instance
- nodes can hardcode Links, VolumesFrom etc to particular instances
