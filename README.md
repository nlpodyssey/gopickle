# GoPickle

GoPickle is a Go library for loading Python's data serialized with `pickle`
and PyTorch module files.

The `pickle` sub-package provides the core functionality for loading data
serialized with Python `pickle` module, from a file, string, or byte sequence.
All _pickle_ protocols from 0 to 5 are supported.

The `pytorch` sub-package implements types and functions for loading
PyTorch module files. Both the _modern_ zip-compressed format and the 
_legacy_ non-tar format are supported. Legacy tar-compressed
files and TorchScript archives are _not_ supported.

## Project Status and Contributions

This project is currently in **alpha** development stage. While we provide
better documentation, tests, and functionalities, above all we'd like to
battle-test this library with real data and models. This would tremendously
help us find tricky bugs, add more built-in types, or change the
library API to make it easier to use.

The simplest and most useful and way to contribute is to try this library
yourself and give us your feedback: report bugs, suggest improvements,
or just tell us your opinion.
And of course, if you feel like it, please go on with your own pull
requests! We can discuss any issue and try to find the best
solution for it together.

## Usage

### Pickle

Simple usage:

```go
import "github.com/nlpodyssey/gopickle/pickle"

// ...

// from file
foo, err := pickle.Load("foo.p") 

// from string
stringDump := "I42\n."
bar, err := pickle.Loads(stringDump)

// ...
```

Advanced/custom usage:

```go
import "github.com/nlpodyssey/gopickle/pickle"

var r io.Reader

// ...

u := pickle.NewUnpickler(r)

// Handle custom classes
u.FindClass = func(module, name string) (interface{}, error) {
    if module == "foo" && name == "Bar" {
        return myFooBarClass, nil
    }
    return nil, fmt.Errorf("class not found :(")
}

// Resolve objects by persistent ID
u.PersistentLoad = func(persistentId interface{}) (interface{}, error) {
    obj := doSomethingWithPersistentId(persistentId)
    return obj, nil
}

// Handle custom pickle extensions
u.GetExtension = func(code int) (interface{}, error) {
    obj := doSomethingToResolveExtension(code)
    return obj, nil
}

// Handle Out-of-band Buffers
// https://docs.python.org/3/library/pickle.html#out-of-band-buffers
u.NextBuffer = func() (interface{}, error) {
    buf := getMyNextBuffer()
    return buf, nil
}

// Low-level function to handle pickle protocol 5 READONLY_BUFFER opcode.
// By default it is completely ignored (sort of no-op); here you have the
// ability to manipulate objects as you need.
u.MakeReadOnly = func(obj interface{}) (interface{}, error) {
    newObj := myReadOnlyTransform(obj)
    return newObj, nil
}

data, err := u.Load()

// ...
```

### PyTorch

The library currently provides a high-level function for loading a module file:

```go
import "github.com/nlpodyssey/gopickle/pytorch"

// ...

myModel, err := pytorch.Load("module.pt")

// ...
```

More features will be provided in the future. 

## How it works

### Pickle

Unlike more traditional data serialization formats, (such as JSON or YAML),
a "pickle" is a _program_ for a so-called _unpickling machine_, also known
as _virtual pickle machine_, or _PM_ for short. A program consists in a
sequence of opcodes which instructs the virtual machine about how to build
arbitrarily complex Python objects. You can learn more  from Python
`pickletools` [module documentation](https://github.com/python/cpython/blob/3.8/Lib/pickletools.py).

Python PM implementation is straightforward, since it can take advantage
of the whole environment provided by a running Python interpreter. For this
Go implementation we want to keep things simple, for example avoiding
dependencies or foreign bindings, yet we want to provide flexibility, and a way
for any user to extend basic functionalities of the library.

This Go unpickling machine implementation makes use of a set of types defined
in `types`.
This sub-package contains Go types representing classes, instances and common
interfaces for some of the most commonly used builtin non-scalar types in 
Python.
We chose to provide only minimal functionalities for each type, for the sole 
purpose of making them easy to be handled by the machine.

Since Python's _pickle_ can dump and load _any_ object, the aforementioned types
are clearly not always sufficient. You can easily handle the loading of any 
missing class by explicitly providing a `FindClass` callback to an `Unpickler`
object. The implementation of your custom classes can be as simple or as
sophisticated as you need. If a certain class is required but is not found,
by default a `GenericClass` is used.
In some circumstances, this is enough to fully load a _pickle_ program, but
on other occasions the pickle program might require a certain class with
specific traits: in this case, the `GenericClass` is not enough and an error
is returned. You should be able to fix this situation providing
a custom class implementation, that jas to reflect the same basic behaviour
you can observe in the original Python implementation.

A similar approach is adopted for other peculiar aspects, such as persistent
objects loading, extensions handling, and a couple of protocol-5 opcodes:
whenever necessary, you can implement custom behaviours providing one or more
callback functions.

Once resolved, all representation of classes and objects are casted to
`interface{}` type; then the machine looks for specific types or
interfaces to be implemented on an object only where strictly necessary. 

The virtual machine closely follows the original implementation
from Python 3.8 - see the [`Unpickler` class](https://github.com/python/cpython/blob/3.8/Lib/pickle.py#L1134). 

### PyTorch

[PyTorch](https://pytorch.org/) machine learning framework allows you to save
and load Python objects which include (or _are_) [Tensors](https://pytorch.org/docs/stable/tensors.html)
or other framework-specific elements. Tensors data are handled by the more
primitive [Storage](https://pytorch.org/docs/stable/storage.html) classes, which
are efficiently serialized as raw sequences of bytes. All the rest is dumped
using `pickle`. When serializing, the programmer can choose any available pickle
protocol, and whether to use zip compression. 

The package `pytorch` implements loading functionalities for data
files serialized with PyTorch (called _modules_). The Go implementation
strictly follows the original Python (and C++) [code](https://github.com/pytorch/pytorch/blob/master/torch/serialization.py#L486).
The `pickle` and `types` packages are used to read some parts of a given file.
Other specific types are implemented in the `pytorch` module itself, most
notably to reflect the content of PyTorch Tensor and Storage objects. 

## License

GoPickle is licensed under a BSD-style license.
See [LICENSE](https://github.com/nlpodyssey/gopickle/blob/master/LICENSE) for
the full license text.
