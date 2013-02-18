function linechart(data) {
    var width = 600, height = 400;
    var margin = {top: 20, right: 20, bottom: 30, left: 60};

    var dataPoints = [];
    for (var series = 0; series < data.length; series++) {
        for (var point = 0; point < data[series].data.length; point++) {
            dataPoints.push(data[series].data[point]);
            data[series].data[point].series_id = data[series].id;
        }
    }

    var wrapperDiv = d3.select('#linechart').
        style('width', (width + margin.left + margin.right) + 'px').
        style('height', (height + margin.top + margin.bottom) + 'px');

    var svg = wrapperDiv.append('svg').
        attr('width', width + margin.left + margin.right).
        attr('height', height + margin.top + margin.bottom).
        attr('class', 'linechart').
        append('g').
        attr('transform', 'translate(' + margin.left + ',' + margin.top + ')');

    // By default, time scale is in the browser's local timezone
    var xScale = d3.time.scale().domain([
        new Date(d3.min(dataPoints, function (d) { return d.timestamp; }) * 1000),
        new Date(d3.max(dataPoints, function (d) { return d.timestamp; }) * 1000)
    ]).range([0, width]).nice();

    var yScale = d3.scale.linear().
        domain([0, d3.max(dataPoints, function (d) { return d.value; })]).
        range([height, 0]).nice();

    var approxYTicks = 4;
    var xAxis = d3.svg.axis().scale(xScale).orient('bottom').ticks(d3.time.days);
    var yAxis = d3.svg.axis().scale(yScale).orient('left').ticks(approxYTicks).tickSubdivide(1);
    var palette = d3.scale.category10();

    svg.append('g').attr('class', 'x axis').
        attr('transform', 'translate(0,' + height + ')').
        call(xAxis);

    var yAxisGroup = svg.append('g').attr('class', 'y axis').call(yAxis);

    yAxisGroup.append('g').attr('class', 'grid').
        selectAll('line').
        data(yScale.ticks(approxYTicks)).enter().
        append('line').
        attr({x1: 0, x2: width, y1: yScale, y2: yScale});

    yAxisGroup.append('text').
        attr('class', 'label').
        attr('transform', 'rotate(-90)').
        attr('dy', '1.2em').
        style('text-anchor', 'end').
        text('seconds');

    var line = d3.svg.line().
        x(function (d) { return xScale(new Date(d.timestamp * 1000)); }).
        y(function (d) { return yScale(d.value); }).
        interpolate('linear');

    var series = svg.selectAll('g.series').data(data).enter().
        append('g').
        attr('class', function (d) { return 'series ' + d.name.replace(/\W+/g, '-'); });

    series.append('path').
        style('stroke', function (d) { return palette(d.id); }).
        attr('d', function (d) { return line(d.data); });

    series.selectAll('circle').data(function (d) { return d.data; }).enter().
        append('circle').
        attr('cx', function (d) { return xScale(new Date(d.timestamp * 1000)); }).
        attr('cy', function (d) { return yScale(d.value); }).
        attr('r', 5).
        style('fill', function(d) { return palette(d.series_id); }).
        style('stroke', function(d) { return palette(d.series_id); });

    var linkSize = 10;

    wrapperDiv.selectAll('a').data(dataPoints.filter(function (d) { return d.url; })).enter().
        append('a').
        attr('href', function (d) { return d.url; }).
        style('width', linkSize + 'px').
        style('height', linkSize + 'px').
        style('left', function (d) { return (xScale(new Date(d.timestamp * 1000)) + margin.left - linkSize / 2) + 'px'; }).
        style('top', function (d) { return (yScale(d.value) + margin.top - linkSize / 2) + 'px'; });
}

jQuery(function () {
    if ($('#linechart').length == 0) return;
    $.getJSON(location.href + '.json', linechart);
});
