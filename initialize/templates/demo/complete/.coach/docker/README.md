# ./coach/docker

An arbitrary place the we put custom Docker builds for the project.

These custom builds are often small overrides of factory images used 
to make alterations, or to add content.

If a node has a "Build: {path}" declaration, then a custom build should
be placed in .coach/{path}/Dockerfile. As a standard we have been using
.coach/docker/{build}/Dockerfile for that "Build: docker/{build}"
 