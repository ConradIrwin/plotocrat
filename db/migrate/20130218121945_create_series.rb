class CreateSeries < ActiveRecord::Migration
  def change
    create_table :series do |t|
      t.string :name, :null => false
      t.timestamps
    end

    add_index :series, :name
  end
end
