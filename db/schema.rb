# encoding: UTF-8
# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended to check this file into your version control system.

ActiveRecord::Schema.define(:version => 20130218121945) do

  create_table "gauge_values", :force => true do |t|
    t.integer  "series_id",  :null => false
    t.float    "value",      :null => false
    t.string   "url"
    t.datetime "created_at", :null => false
  end

  add_index "gauge_values", ["series_id"], :name => "index_gauge_values_on_series_id"

  create_table "plots", :force => true do |t|
    t.text     "slug",       :null => false
    t.text     "data",       :null => false
    t.text     "title"
    t.datetime "created_at", :null => false
    t.datetime "updated_at", :null => false
  end

  add_index "plots", ["slug"], :name => "index_plots_on_slug"

  create_table "series", :force => true do |t|
    t.string   "name",       :null => false
    t.datetime "created_at", :null => false
    t.datetime "updated_at", :null => false
  end

  add_index "series", ["name"], :name => "index_series_on_name"

end
