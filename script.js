// Include this script in an HTML file
// Ensure you include D3.js in your project by adding:
// <script src="https://d3js.org/d3.v7.min.js"></script>

const data = [
  { day: 1, values: [1540, 614, 124, 151, 119, 440, 81, 61, 192, 62] },
  { day: 2, values: [641, 443, 132, 183, 350, 136, 104, 226, 171, 150] },
  { day: 3, values: [410, 263, 482, 268, 433, 165, 326, 178, 421, 80] },
  { day: 4, values: [240, 744, 92, 617, 184, 444, 342, 101, 211, 95] },
  { day: 5, values: [367, 484, 147, 289, 667, 195, 237, 241, 768, 72] },
  { day: 6, values: [653, 158, 263, 613, 356, 140, 154, 77, 170, 156] },
  { day: 7, values: [1281, 636, 647, 765, 910, 451, 108, 463, 568, 69] },
  { day: 8, values: [365, 581, 217, 311, 333, 264, 526, 313, 318, 108] },
  { day: 9, values: [481, 648, 333, 665, 344, 211, 349, 468, 161, 151] },
  { day: 10, values: [212, 713, 641, 413, 1253, 240, 281, 405, 907, 60] },
  { day: 11, values: [484, 1031, 286, 456, 438, 360, 307, 611, 270, 99] },
  { day: 12, values: [393,470,186,666,979,352,370,282,602,500] },
  { day: 13, values: [557,491,462,1259,658,441,309,423,432,116] },
  { day: 14, values: [642,917,591,496,884,450,351,488,545,349] },
  { day: 15, values: [591,345,219,2923,1144,293,392,900,231,1059] },
  { day: 16, values: [460,351,536,1088,1514,642,739,1586,501,105] },
  { day: 17, values: [239,607,457,2449,1309,384,379,1384,515,1340] },
  { day: 18, values: [608,389,933,638,2335,276,1280,346,493,74] },
  { day: 19, values: [1747,1248,455,1248,641,752,1705,2062,994,55] },
  { day: 20, values: [572,402,558,1225,1370,2016,525,576,1360,267] },
  { day: 21, values: [886,1165,984,1041,936,519,575,467,1410,2099] },
  { day: 22, values: [2212,1415,516,1288,1665,650,919,2360,914,118] },
  { day: 23, values: [493,1356,1382,1620,561,1035,2102,780,793,67] },
  { day: 24, values: [545,791,496,2312,1173,480,1690,744,1266,1585] },
  { day: 25, values: [361,329,304,213,882,197,287,222,289,90] },
];

const years = [2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024];

// Set dimensions
const margin = { top: 20, right: 80, bottom: 50, left: 80 };
const width = 2000 - margin.left - margin.right;
const height = 1000 - margin.top - margin.bottom;

// Create SVG container
const svg = d3.select("body")
  .append("svg")
  .attr("width", width + margin.left + margin.right)
  .attr("height", height + margin.top + margin.bottom)
  .attr("style", "max-width: 100%; height: auto; font: 20px sans-serif;")
  .append("g")
  .attr("transform", `translate(${margin.left},${margin.top})`);

// Define scales
const xScale = d3.scaleLinear()
  .domain([data[0].day, data[data.length-1].day]) // Days
  .range([0, width]);

const yScale = d3.scaleLinear()
  .domain([0, d3.max(data, d => d3.max(d.values))])
  .nice()
  .range([height, 0]);

// Define color scale
const colorScale = d3.scaleOrdinal(d3.schemeCategory10)
  .domain(years)

// Define axes
const xAxis = d3.axisBottom(xScale).ticks(11).tickFormat(d => `Day ${d}`);
const yAxis = d3.axisLeft(yScale);

// Add grid lines
svg.append("g")
  .attr("class", "grid")
  .attr("transform", `translate(0,${height})`)
  .attr("opacity", 0.10)
  .call(d3.axisBottom(xScale)
    .tickSize(-height)
    .tickFormat(""));

svg.append("g")
  .attr("class", "grid")
  .attr("opacity", 0.10)
  .call(d3.axisLeft(yScale)
    .tickSize(-width)
    .tickFormat(""));

// Add axes to the SVG
svg.append("g")
  .attr("transform", `translate(0,${height})`)
  .call(xAxis);

svg.append("g")
  .call(yAxis);

// Add axis labels
// svg.append("text")
//   .attr("x", width / 2)
//   .attr("y", height + margin.bottom - 10)
//   .attr("text-anchor", "middle")
//   .text("Day");

svg.append("text")
  .attr("x", -height / 2)
  .attr("y", -margin.left + 30)
  .attr("text-anchor", "middle")
  .attr("transform", "rotate(-90)")
  .attr("style", "max-width: 100%; height: auto; font: 20px sans-serif;")
  .text("Average Seconds of First Gold Stars");

// Line generator
const line = d3.line()
  .x((d, i) => xScale(i + 1))
  .y(d => yScale(d));

// Draw lines and points for each year
years.forEach((year, yearIndex) => {
  const yearData = data.map(d => d.values[yearIndex]);

  // Draw line
  svg.append("path")
    .datum(yearData)
    .attr("fill", "none")
    .attr("stroke", colorScale(year))
    .attr("stroke-width", 2)
    .attr("d", line);

  // Draw dots
  svg.selectAll(`.dot-${year}`)
    .data(yearData)
    .enter()
    .append("circle")
    .attr("class", `dot-${year}`)
    .attr("cx", (d, i) => xScale(i + 1))
    .attr("cy", d => yScale(d))
    .attr("r", 4)
    .attr("fill", colorScale(year));

  // Add legend
  svg.append("text")
    .attr("x", 50)
    .attr("y", 20 + yearIndex * 20)
    .attr("fill", colorScale(year))
    .text(year);
});

// Add download button
const downloadButton = d3.select("body").append("button")
  .text("Download Chart as Image")
  .on("click", () => {
    const serializer = new XMLSerializer();
    let svgString = serializer.serializeToString(svg.node());
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    const img = new Image();

    canvas.width = width + margin.left + margin.right;
    canvas.height = height + margin.top + margin.bottom;

    img.onload = () => {
      ctx.drawImage(img, 0, 0);
      const a = document.createElement("a");
      a.download = "chart.png";
      a.href = canvas.toDataURL("image/png");
      a.click();
    };

    img.onerror = (err) => {
      console.error(err)
    }

    svgString = `<svg xmlns="http://www.w3.org/2000/svg">` + svgString + `</svg>`
    img.src = 'data:image/svg+xml;base64,' + btoa(svgString);
  });
