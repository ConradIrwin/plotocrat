class PlotsController < ApplicationController
  def index
  end

  def upload
    title, upload = params.detect do |k, v|
      ActionDispatch::Http::UploadedFile === v
    end

    unless upload
      return render :status => 401, :error => "Usage: cat data | curl plotocrat.com -F 'X axis title'=@-"
    end

    plot = Plot.new(:data => upload.read, :title => title).tap(&:save!)
    render :text => url_for(:action => :view, :slug => plot.slug) + "\n"
  rescue ActiveRecord::RecordInvalid => e
    render :status => 401, :text => url_for(:action => :error, :error => e.message)
  end

  def view
    @plot = Plot.where(:slug => params[:slug]).first
  end

  def error
    render :text => params[:error]
  end

  private

  def render(args={})
    if error = args.delete(:error)
      super :text => url_for(:action => :error, :error => e.message)
    else
      super
    end
  end
end
