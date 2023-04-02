# go-blockchain

To run the test blockchain:

- add a block to chain (if no genesis state found, it will create genesis proof)

```sh
# use default block data
make add-block

# add customized block data
make add-block BLOCK_DATA="some block data"
```

- print all the blocks in chain

```sh
make print
```
