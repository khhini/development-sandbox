const { default: Automizer, modify, DxaToCm } = require("pptx-automizer/dist");


const run = async () => {

  const automizer = new Automizer({
    templateDir: './'
  });

  let pres = automizer
    .loadRoot("reports.pptx") // Load reports.pptx as Root Template 
    .load("reports.pptx"); // Load reports.pptx as working presentations

  // Add slide number 1 from root templates to working presentations.
  pres.addSlide("reports.pptx", 1, (slide) => {
    slide.getAllElements().then((elements) => {
      elements.forEach((element) => {
        console.log(element.name)
      })
    })

    // Modify Text of Element named "Subtitle 2"
    slide.getElement("Subtitle 2").then((element) => {
      console.log(element)
      slide.modifyElement(element.name, modify.setText("Test Modifying a String"));
    })
  })

  // Add slide number 1 from root templates to working presentations.
  pres.addSlide("reports.pptx", 2, (slide) => {
    slide.getAllElements().then((elements) => {
      elements.forEach((element) => {
        if (element.name.includes("Picture")) {
          console.log(element.position)


          // Generate new chart elements
          const dataChartAreaLine = [
            {
              name: 'Actual Sales',
              labels: ['Jan', 'Feb', 'Mar'],
              values: [1500, 4600, 5156],
            },
            {
              name: 'Projected Sales',
              labels: ['Jan', 'Feb', 'Mar'],
              values: [1000, 2600, 3456],
            },
          ];

          slide.generate((pSlide, pptxGenJS) => {
            pSlide.addChart(pptxGenJS.ChartType.line, dataChartAreaLine, {
              x: DxaToCm(element.position.x),
              y: DxaToCm(element.position.y),
              w: DxaToCm(element.position.cx),
              h: DxaToCm(element.position.cy)
            })
          })

          // Delete Elements
          slide.removeElement(element.name)
        }
      })


    })
  })

  // Write Presentations to results.pptx files
  pres.write("results.pptx").then((summary) => {
    console.log(summary)
  })

}

run().catch((error) => {
  console.error(error);
})
