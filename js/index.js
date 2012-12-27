(function () {
    $(function () {
        // Ideas from gnuplot's quantize_normal_tic
        //
        // returns between 5 and 10 tics
        function axisTics(data) {
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
        }

        var height = 400, width = 800;

        function raw(data) {

            data.sort(d3.ascending);
            var axis = axisTics(data);
            axis[0] = Math.max(axis[0], 1);

            var x = d3.scale.linear().domain([0, data.length - 1]).range([40, width]);
            var px = d3.scale.linear().domain([0, 100]).range([40, width]);
            var y = d3.scale.linear().domain([d3.min(axis), d3.max(axis)]).range([height - 20, 20]);

            var viz = d3.select('#chart').style('width', width).style('height', height)
                        .append('svg:svg').attr('width', width).attr('height', height).attr('class', 'viz')
                        .append('svg:g');

            viz.selectAll('path.line').data([data]).enter().append('svg:path')
                .attr('d', d3.svg.line().x(function (d, i) { return x(i); }).y(y));

            viz.selectAll('.tickx').data([0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100]).enter().append('svg:line')
                .attr('y1', y(axis[0])).attr('y2', y(axis[0]) - 3).attr('x1', px).attr('x2', px);

            viz.selectAll('.ticky').data(axis).enter().append('svg:line')
                .attr('x1', x(0)).attr('x2', x(0) + 3).attr('y1', y).attr('y2', y);

            viz.selectAll('.labely').data(axis).enter().append('svg:text')
                .text(function (d) { return d; }).attr('x', 0).attr('y', function (d) { return y(d) + 5; });

        }

        // http://www.umiacs.umd.edu/labs/cvl/pirl/vikas/Software/optimal_bw/optimal_bw_code.htm
        function cdf(data) {
            var x = d3.scale.log().domain([Math.max(1, d3.min(data)), d3.max(data)]).range([40, width]);
            var y = d3.scale.linear().domain([1, 0]).range([40, height]);
            var kde = science.stats.kde().sample(data.map(function (x) { return Math.log(x); }))(d3.range(Math.log(data[0]), Math.log(data[data.length - 1]), 0.05));
            var xk = d3.scale.linear().domain([Math.log(data[0]), Math.log(data[data.length - 1])]).range([40, width]);
            var yk = d3.scale.linear().domain([d3.max(kde.map(function (x) { return x[1]; })), d3.min(kde.map(function (x) { return x[1]; }))]).range([40, height]);

            var viz = d3.select('#cdf').style('width', width).style('height', height)
                        .append('svg:svg').attr('width', width).attr('height', height).attr('class', 'viz')
                        .append('svg:g');

            var cumulative = [], i = 0;

            while (i < data.length) {
                while (data[i + 1] === data[i]) {
                    i += 1;
                }
                cumulative.push([data[i], (i + 1) / data.length]);
                i += 1;
            }

            viz.selectAll('path.pdf').data([kde]).enter().append('svg:path')
                .attr('class', 'pdf')
                .attr('d', d3.svg.line().x(function (d) {
                    return xk(d[0]);
                }).y(function (d) {
                    return yk(d[1]);
                }));

            viz.selectAll('path.cdf').data([cumulative]).enter().append('svg:path')
                .attr('class', 'cdf')
                .attr('d', d3.svg.line().x(function (d) {
                    return x(d[0]);
                }).y(function (d) { return y(d[1]); }));


        }

        window.emails = window.emails.filter(function (x) { return x !== 0; });

        window.normal = window.normal.filter(function (x) { return x !== 0; }).sort(d3.ascending);

        raw(window.emails.map(Number));
        cdf(window.emails.map(Number));
    });
}).call(this);
