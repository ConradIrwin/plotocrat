(->
  $ ->
    cdf = (data) ->
      kde = science.stats.distribution.kde().sample(data).resolution(200)
      x = undefined

      isIntegral = d3.max(data, (x) -> x % 1) == 0
      if kde.feelsLogarithmic()
        d2 = []
        # TODO: figure out the correct solution...
        data.forEach (x) -> d2.push(x) unless x is 0
        data = d2
        kde = science.stats.distribution.kde().sample(data).resolution(200)
        x = d3.scale.log()
      else
        x = d3.scale.linear()

      kde.scale x.copy()
      x.clamp(true).domain([data[0], data[data.length - 1]]).range([40, width]).nice()
      axisTicks = x.ticks(5)

      y = d3.scale.linear().domain([1, 0]).range([40, height])
      yk = d3.scale.linear().domain([kde.max(), 0]).range([40, height])

      label = d3.select("#cdf").append("div").html(title + " = <br/>cumulative =<br/>density =")
      title = d3.select('#cdf').attr('data-title')

      svg = d3.select("#cdf")
        .style("width", width)
        .style("height", height)
        .append("svg:svg")
        .attr("width", 40 + width)
        .attr("height", height + 40)
        .attr("class", "viz")
      viz = svg
        .append("svg:g")

      viz.selectAll("line.ydivisions.strong").data([0, 1]).enter()
        .append("svg:line")
        .attr
          class: "ydivisions strong"
          x1: x(axisTicks[0])
          x2: x(axisTicks[axisTicks.length - 1])
          y1: y
          y2: y

      viz.selectAll("line.ydivisions.weak").data(d3.range(0, 1, 0.25)).enter()
        .append("svg:line")
        .attr
          class: "ydivisions weak"
          x1: x(axisTicks[0])
          x2: x(axisTicks[axisTicks.length - 1])
          y1: y
          y2: y

      viz.selectAll("line.xdivisions").data(axisTicks).enter()
        .append("svg:line")
        .attr
          class: (d, i) ->
            if i is 0 or d.toString().match(/^10*$/) or i is axisTicks.length - 1
              "xdivisions strong"
            else
              "xdivisions weak"
          y1: y(0)
          y2: y(1)
          x1: x
          x2: x

      viz.selectAll("text.xticklabels").data(axisTicks).enter()
        .append("svg:text").attr("class", "xticklabels")
        .text((d, i) ->
          raw = x.tickFormat(5)(d)

          raw.replace(/^1e\+([0-6])$/, (m, i) ->
            Math.pow(10, i)
          )

        ).attr
          x: x
          y: height + 20

      viz.selectAll("text.xtitle").data([0]).enter()
        .append("svg:text").attr("class", "xtitle")
        .text(title)
        .attr
          x: 40 + width / 2
          y: height + 60

      viz.selectAll("text.ytitle").data([0]).enter()
        .append("svg:text").attr("class", "ytitle")
        .text("Probability")
        .attr
          'text-anchor': 'middle'
          transform: "rotate(-90)"
          x: - (40 + height) / 2
          y: 20

      viz.selectAll("text.y2title").data([0]).enter()
        .append("svg:text")
        .attr("class", "y2title")
        .text("Cumulative Probability")
        .attr
          transform: "rotate(90)"
          x: (height + 40) / 2
          y: - (width + 20)

      viz.selectAll("path.pdf").data([kde.pdf()]).enter()
        .append("svg:path")
        .attr
          class: "pdf"
          d: d3.svg.line().x((d) -> x d[0]).y((d) -> yk d[1])

      viz.selectAll("path.cdf").data([kde.cdf()]).enter()
        .append("svg:path")
        .attr
          class: "cdf"
          d: d3.svg.line().x((d) -> x d[0]).y((d) -> y d[1])

      viz.selectAll("circle.mode").data([kde.mode()]).enter()
        .append("svg:circle")
        .attr
          class: "mode"
          r: 4
          cx: x
          cy: (d) -> yk kde(d)

      viz.selectAll("circle.expectation").data([kde.expectation()]).enter()
        .append("svg:circle")
        .attr
          class: "expectation"
          r: 4
          cx: x
          cy: -> y 0.25

      viz.selectAll("circle.percentile").data(kde.qf()).enter()
        .append("svg:circle")
        .attr
          class: "percentile"
          r: 4
          cx: (d) -> x d[0]
          cy: (d) -> y d[1]


      show_position = (raw_pos) ->
        value = x.invert(raw_pos)
        cumulative = kde.inverseQuantile(value)
        density = kde(value)

        if value < axisTicks[0]
          value = axisTicks[0]
        else if value > axisTicks[axisTicks.length - 1]
          value = axisTicks[axisTicks.length - 1]

        line = viz.selectAll("line.fugi").data([value])
        line.enter().append("svg:line").attr
          class: "fugi"
          y1: y(0)
          y2: y(1)
        line.attr
          x1: x
          x2: x


        cd = [yk(density), y(cumulative)]

        circles = viz.selectAll("circle.fugi").data(cd)
        circles.enter().append("svg:circle").attr
          r: 5
          class: 'fugi'
        circles.attr
          cx: x(value)
          cy: (d) -> d

        value_s = title + " = " + if isIntegral then value.toFixed() else value.toPrecision(4)
        value_s = value_s.replace(/\((.*)\) *= ([0-9e\-\+\.]*)$/, (m, u, v) -> "= " + v + " " + u)
        cumulative_s = "cumulative = " + (cumulative * 100).toPrecision(4) + "%"
        absolute_s = "density = " + density.toPrecision(3)

        label.html([value_s, cumulative_s, absolute_s].join("<br/>"))

      toggle_fugi = (visible) ->
        viz.selectAll(".fugi").style
          display: if visible then 'block' else 'none'

      svg.on('mousemove', () ->
        show_position d3.mouse(svg.node())[0]
      ).on('mouseover', ->
        toggle_fugi true
      ).on('mouseout', ->
        toggle_fugi false
      )

      is_multitouch = false
      is_scrolling = null
      start = []
      svg.on('touchmove', (e) ->
        touches = d3.touches svg.node()

        # ensure that we don't insert a line after the user releases
        # one finger from zoom-pinching
        is_multitouch ||= touches.length > 1

        # If the user is moving one finger up or down, they probably want to scroll
        if is_scrolling == null
          is_scrolling = Math.abs(touches[0][0] - start[0][0]) < Math.abs(touches[0][1] - start[0][1])

        if is_multitouch || is_scrolling
          toggle_fugi false
          return

        toggle_fugi true
        show_position touches[0][0]
        # prevent scrolling
        d3.event.preventDefault()

      ).on('touchstart', ->
        start = d3.touches svg.node()
        is_scrolling = null
        is_multitouch = false
      ).on('touchend', ->
        toggle_fugi false
      )


    height = 400
    width = 600
    data = []
    $('textarea#data').val().split("\n").forEach((x) ->
      return if x.match(/^(#|\s*$)/)
      n = Number(x)
      # TODO warn about unparseable lines?
      data.push(n) unless isNaN(n)
    )

    data.sort(d3.ascending)

    cdf data

).call this
