class ApplicationController < ActionController::Base
  protect_from_forgery

  private
  def render(args={})
    if Hash === args && args[:text]
      response.headers['Content-Type'] = 'text/plain; charset=utf-8'
    end
    super
  end
end
