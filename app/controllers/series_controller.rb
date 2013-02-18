class SeriesController < ApplicationController
  # GET /series
  # GET /series.json
  def index
    @series = Series.all

    respond_to do |format|
      format.html # index.html.erb
      format.json { render json: @series }
    end
  end

  # GET /series/1
  # GET /series/1.json
  def show
    @series = Series.where(:id => params[:id].split(','))
    @data = @series.flat_map(&:data)

    respond_to do |format|
      format.html # show.html.erb
      format.json { render json: @series }
    end
  end

  # POST /series
  # POST /series.json
  def create
    @data_point = GaugeValue.new_from_params(params)

    respond_to do |format|
      if @data_point.save
        format.html { redirect_to @data_point.series }
        format.json { render status: :created, location: @data_point.series }
      else
        format.json { render json: @data_point.errors, status: :unprocessable_entity }
      end
    end
  end
end
