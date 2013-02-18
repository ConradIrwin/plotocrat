class CreateGaugeValues < ActiveRecord::Migration
  def change
    create_table :gauge_values do |t|
      t.integer  :series_id,  :null => false
      t.float    :value,      :null => false
      t.string   :url,        :null => true
      t.datetime :created_at, :null => false
    end

    add_index :gauge_values, :series_id
  end
end
