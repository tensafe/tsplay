cleanup = db_execute({
    connection = "reporting",
    sql = "DELETE FROM public.tutorial_orders WHERE order_id = $1",
    args = {"T1001"}
})

insert_result = db_insert({
    connection = "reporting",
    table = "public.tutorial_orders",
    columns = {"order_id", "status", "amount"},
    row = {
        order_id = "T1001",
        status = "ready",
        amount = 99.5
    }
})

row = db_query_one({
    connection = "reporting",
    sql = "SELECT order_id, status, amount FROM public.tutorial_orders WHERE order_id = $1",
    args = {"T1001"}
})

write_json("artifacts/tutorials/07-db-postgres-basics-lua.json", {
    lesson = "07",
    mode = "lua",
    cleanup = cleanup,
    insert_result = insert_result,
    row = row
})

print("inserted order:", tostring(row.order_id))
print("order status:", tostring(row.status))
print("wrote artifacts/tutorials/07-db-postgres-basics-lua.json")
