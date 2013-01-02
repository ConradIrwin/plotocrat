class Plot < ActiveRecord::Base
  attr_accessible :slug, :title, :data
  validates :data, :presence => true
  before_save   :ensure_crlf
  before_create :set_defaults

  def ensure_crlf
    self.data = self.data.gsub(/\r\n|\r|\n/, "\r\n")
  end

  def set_defaults
    self.slug = Digest::SHA1.hexdigest(data)[0...20] if data
  end
end
