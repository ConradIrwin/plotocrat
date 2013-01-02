class UploadsController < ActionController::Metal

  def upload
    plot = Plot.new(:data => request.raw_post, :title => params[:title]).tap(&:save!)
    puts request.raw_post.inspect
    render :text => url_for(:action => :view, :slug => plot.slug)
  rescue ActiveRecord::RecordInvalid => e
    render :status => 401, :text => url_for(:action => :error, :error => e.message)
  end

end
