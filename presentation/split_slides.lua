function Pandoc(doc)
  local slides = {}
  local current_slide = {}

  for _, block in ipairs(doc.blocks) do
    if block.t == "HorizontalRule" then
      if #current_slide > 0 then
        table.insert(slides, pandoc.Div(current_slide, {class="section slide"}))
        current_slide = {}
      end
    else
      table.insert(current_slide, block)
    end
  end

  if #current_slide > 0 then
    table.insert(slides, pandoc.Div(current_slide, {class="section slide"}))
  end

  return pandoc.Pandoc(slides, doc.meta)
end
