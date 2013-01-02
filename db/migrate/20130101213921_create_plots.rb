class CreatePlots < ActiveRecord::Migration
  def change
    create_table :plots do |t|
      t.text :slug,  :null => false
      t.text :data,  :null => false
      t.text :title, :null => true

      t.timestamps
    end

    add_index :plots, :slug
  end
end
