# Chatcorr

This (hopefully) implements the [Chatterjee Correlation](https://arxiv.org/pdf/1909.10140.pdf).
The implementation here follows the paper, but has not been validated against any other implementations (yet).
This is also an exercise in using Go generics.
No attempt was made to tune this for performance.

This code supplies 4 pairs of functions for computing correlations of Points, which are generic in T and U.

Because the Chatterjee Correlation chooses randomly where there are ties in the X coordinates, each pair
of functions includes one that creates its own random number generator, and one that takes an RNG as a
parameter to allow repeatable results.  In all cases, the function sorts its input slice of points by
their Y coordinates.

The first function `CCF64([]Point[float64, float64]) float64` is monomorphic, because float64 coordinates are a likely case.

The second function `CC[T, U Lessable](v []Point[T, U]) float64` works for any coordinates where the two types are `Lessable`, i.e., are one of the primitive types with a comparison operator.

The third function `CCFn[T any](v []Point[T, T], compare func (a,b T) int) float64` works for any coordinates that share the same type, and orders them using the supplied function.

The fourth function `CCMixed[T,U any](v []Point[T, U], compareT func (a,b T) int, compareU func (a,b U) int) float64` is fully general and works for any pair of types that can be compared by the two supplied functions.
