Factors and Caching challenge
=============================

Usage:

    ruby main.rb [arguments]

    arguments:
        space-separated list of base-10 integers. An error will be printed to stderr for each argument that could not be parsed, although the program will not be interrupted.

    example:
        ruby main.rb 10 5 2 20


1.  What if you were to cache the calculation, for example in the file system.  What would an example implementation of the cache look like?  By cache I mean, given an array input, skip the calculation of the output if you have already calculated the output at least once already.

If I were to cache this program, I would use a key-value store -- probably redis, but, depending on situation, i would consider mysql, memcached, and raw files.

To generate the keys, I would use the array of parsed integers -- intentionally discarding invalid integers -- and order them in ascending order. This is to ensure that the key is generated from the same data regardless of the order of the input.

I would then take the ordered list of integers, convert each to a string, and concatenate them together with a non-numeric separator, such as comma. Then, I would pass the concatenated string through a 1-way deterministic hashing function -- sha256, for example. This hash would be my key (filename if i was storing on disk). The purpose of this is to decrease the size of the key from potentially massive to a uniform size. This is safe as long as the hashing function is deterministic.

Then, to generate the value, i would pass the calculated output through a serialization function -- json by default, although a language-oriented serializer might make sense depending on the implementation and performance requirements (such as php serialize) 

Then, to retrieve the output, one only needs to calculate the key again. As the key is calculated predictably -- by ordering the numbers predictably, and by using a deterministic hashing function, no collisions are guaranteed, and the cache will not miss due to a miscalculation.

2. What is the performance of your caching implementation? Is there any way to make it more performant.

The performance of the caching implementation is probably detrimental to this specific problem. The time required to generate a reliable key and communicate outward (either to disk or to a service/redis) is almost always going to induce more overhead/disk wait than simply recalculating the value.

If performance was absolutely required here, the simplest gain is optimization passes over the code, and, if that is not sufficient, to rewrite this solution again in a more performant language, and call the code from ruby.

I wrote what I believe to be the most performant solution. I'm not certain that I am correct, particularly because I do not understand all of the semantics of ruby. If I were to attempt to optimize this further, I would begin by looking into the performance of hashmaps in ruby.

I use hashmaps for the root level of the output variable (output\[current_int\]), and for the individual factors (output\[current_int\]\[factor\]). This is used for ease of lookup, and for free deduplication on the inner layer. It is possible to implement this with just plain arrays, although with much more effort. I do not know if this will be faster or not, as I am uncertain of the overhead induced by hashmaps.

3.  What if you wanted to reverse the functionality.  What if you wanted to output each integer and all the other integers in the array that is the first integer is a factor of I.E:

This is a derived set. The information already exists within the current output. No change to caching is required. An implementation of this transformation would look like:

```ruby
def transform(output)
    new_output = {};
    output.each{ |current_int, factors|
        factors.each{ |factor, discard| 
            new_output.fetch(factor, {})[current_int] = true;
        }
    }
    return new_output
end
```
