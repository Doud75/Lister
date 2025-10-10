UPDATE setlist_items
SET notes = i.script
    FROM interludes i
WHERE setlist_items.item_type = 'interlude'
  AND setlist_items.interlude_id = i.id;