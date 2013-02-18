d3.svg.basicPlot = (container, width, height, x, y) ->
  viz = container
      .style("width", width)
      .style("height", height)
      .append("svg:svg")
      .attr("width", 40 + width)
      .attr("height", height + 40)
      .attr("class", "viz")

  axisTicks = x.ticks(5)

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

  viz
    .append("svg:text").attr("class", "ytitle")
    .text("Probability")
    .attr
      'text-anchor': 'middle'
      transform: "rotate(-90)"
      x: - (40 + height) / 2
      y: 20

  viz
    .append("svg:text")
    .attr("class", "y2title")
    .text("Cumulative Probability")
    .attr
      transform: "rotate(90)"
      x: (height + 40) / 2
      y: - (width + 20)


  viz.onHorizontalInteraction = (f) ->
    f_scaled = (raw_pos) ->
      value = x.invert(raw_pos)

      if value < axisTicks[0]
        value = axisTicks[0]
      else if value > axisTicks[axisTicks.length - 1]
        value = axisTicks[axisTicks.length - 1]

      f value

    viz.on('mousemove', () ->
      f_scaled d3.mouse(viz.node())[0]
    ).on('mouseout', ->
      f false
    )

    is_multitouch = false
    is_scrolling = null
    start = []
    viz.on('touchmove', (e) ->
      touches = d3.touches viz.node()

      # ensure that we don't insert a line after the user releases
      # one finger from zoom-pinching
      is_multitouch ||= touches.length > 1

      # If the user is moving one finger up or down, they probably want to scroll
      if is_scrolling == null
        is_scrolling = Math.abs(touches[0][0] - start[0][0]) < Math.abs(touches[0][1] - start[0][1])

      if is_multitouch || is_scrolling
        f false
        return

      f_scaled touches[0][0]
      # prevent scrolling
      d3.event.preventDefault()

    ).on('touchstart', ->
      start = d3.touches viz.node()
      is_scrolling = null
      is_multitouch = false
    ).on('touchend', ->
      f false
    )



  viz

