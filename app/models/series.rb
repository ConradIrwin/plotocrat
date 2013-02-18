class Series < ActiveRecord::Base
  has_many :gauge_values
  attr_accessible :name
  validates :name, presence: true

  def data
    gauge_values.map do |value|
      {timestamp: value.created_at.to_i, value: value.value, url: value.url}
    end
  end

  def as_json(options=nil)
    super.merge(data: data)
  end
end
