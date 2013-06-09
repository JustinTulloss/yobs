class CreateUsersAndTransactions < ActiveRecord::Migration
  def self.up
    create_table :users do |t|
      t.column :facebook_id, :integer, :null => false
    end
    create_table :transactions do |t|
      t.column :owner_id, :integer, :null => false
      t.column :amount, :integer, :null => false, :default => 0
      t.column :description, :text
    end
  end

  def self.down
    drop_table :users
    drop_table :transactions
  end
end
