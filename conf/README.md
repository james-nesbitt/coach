# Conf

The configuration system provides a struct that contains key configuration for 
a project, including important paths, and string value tokens.  A project conf
struct can be manually created, but is typically create from a combination of
the project root path, and any conf.yml in a coach configuration folder.

Project Name: used as a part of machine names for images and containers in a project
Paths: used as paths and as tokens for a project
Tokens: a string map used for string substitution
