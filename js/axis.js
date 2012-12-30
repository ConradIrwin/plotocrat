science.axis = {};
science.axis.logTicks = function (data) {
    var min = data[0],
        max = data[data.length - 1],

        axis_min = Math.floor(Math.log(min) / Math.log(10)),
        axis_max = Math.ceil(Math.log(max) / Math.log(10)),
        pattern,
        ret = [];

    if (axis_max - axis_min < 1) {
        pattern = [1, 2, 3, 4, 5, 6, 7, 8, 9];

    } else if (axis_max - axis_min < 4) {
        pattern = [1, 2, 4, 6, 8];

    } else if (axis_max - axis_min < 8) {
        pattern = [1, 2.5, 5, 7.5];

    } else {
        pattern = [1, 5];

    }

    for (var i = axis_min; i < axis_max; i += 1) {
        ret = ret.concat(pattern.map(function (t) {
            return Math.pow(10, i) * t;
        }));
    }

    ret.push(Math.pow(10, axis_max));

    return ret;
};


// Ideas from gnuplot's quantize_normal_tic
//
// returns between 5 and 10 tics
science.axis.linearTicks = function (data) {
    var min = data[0],
        max = data[data.length - 1],
        range = max - min,
        size,
        axis_min, axis_max,
        ret = [];
    // Ensure min and max don't crash into the axis.
    min -= range * 0.05;
    max += range * 0.05;
    size = Math.pow(10, Math.floor(Math.log(range) / Math.log(10)));

    // By construction 1 <= range / size < 10 as what we're doing is:
    //     range / 10 ** floor(log10(range))
    // Having a graph with fewer than 5 tics is not great, so we scale the axis if
    // this would happen. We use factors of 2 and 5 because humans are used to
    // them.
    if (range / size < 2) {
        // [0, 2, 4, 6]
        size *= 0.2;

    } else if (range / size < 5) {
        // [0, 5, 10, 15]
        size *= 0.5;

    } else {
        // [0, 1, 2, 3]
        size *= 1;
    }

    axis_min = min - (min % size);
    axis_max = max + size - (max % size);

    for (var i = axis_min; i <= axis_max; i += size) {
        ret.push(i);
    }

    return ret;

};
