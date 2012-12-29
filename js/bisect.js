/* Binary search within an array.
 *
 * Returns the first index at which 'f' returns false.
 *
 * e.g. to find the index of 3 in array:
 *   science.bisect(array, function (x) { return x < 3;})
 *
 * Taken from http://hg.python.org/cpython/file/2.7/Lib/bisect.py
 */
science.bisect = function (array, f) {
    var lo = 0, hi = array.length, mid;

    while (lo < hi) {
        mid = Math.floor((lo + hi) / 2);

        if (f(array[mid])) {
            lo = mid + 1;
        } else {
            hi = mid;
        }
    }

    return lo;
};
