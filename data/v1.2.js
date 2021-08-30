db.form_data.updateMany({
    'is_delete':{ $exists: false }
},{
    $set:{'is_delete':false}
})

db.getCollection("data_log").drop();
db.createCollection("data_log");

db.getCollection("data_log").createIndex({
    "form_id":NumberInt("1")
},{
    name:"form"
});

db.getCollection("data_log").createIndex({
    "action":NumberInt("1")
},{
    name:"action"
})
