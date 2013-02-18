class GaugeValue < ActiveRecord::Base
  belongs_to :series
  attr_accessible :series_id, :value, :url
  validates :series_id, :value, presence: true
  validates_associated :series

  def self.new_from_params(params)
    data_point = new(value: params[:value], url: params[:url])
    data_point.series = Series.where(name: params[:series]).first_or_create if params[:series].present?
    data_point.series = Series.find(params[:series_id]) if params[:series_id].present?
    data_point.created_at = Time.at(params[:timestamp].to_i) if params[:timestamp]
    data_point
  end
end
