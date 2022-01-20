# Graffe

Graffe is a directed graph configuration language meant to be written by humans.
The syntax is very similar to TOML, but it is much simpler.

```toml
# this is a comment
[root]                    # This defines the start of a node

# nodes can have a series of KeyValue pairs

leafs = {one, two, three} # arrays are supported
is_root = true;           # semicolons are optional
num_children = 3          # numbers have a few flavors
pi = 3.14159
hex = 0xFFFF
bin = 0b010101

[one]
descr = "leftmost leaf of root"
[two]
descr = "middle leaf of root"
[three]
descr = "rightmost leaf of root"
```

The main datastructure is the Graph Node, but each Node has an underlying Associative Array.
The parser is multipass, declarations of Nodes can happen in any order.
Values can only refer to nodes, not to keys. Keys live in the Node's namespace, Nodes live in global namespace.
