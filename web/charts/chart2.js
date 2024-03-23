function MoreChartOptions(){} 
var ChartData2 = {
	labels : ["2020","2021","2022","2023",],
	datasets : [
		{
		data : [50000,75000,90000,120000,],
		backgroundColor :'#a478f0',
		borderColor : 'rgba(136,136,136,0)',
		pointBackgroundColor:'#3498db',
		pointBorderColor : '#fff',
		label:"Проектов, использующих ML"},

]
	};
ChartOptions2= {
responsive:false,
	layout:{padding:{top:12,left:12,bottom:12,},},
	 scales: {
	xAxes:[{
gridLines:{color:'rgba(0,0,0,0)',borderDash:[],},
}],

	yAxes:[{
gridLines:{color:'rgba(0,0,0,0)',borderDash:[],},
}],
 },plugins:{
datalabels:{display:true,
	color:'rgba(0,0,0,0)',
	font:{
		style:' bold',},},
},
legend:{
	labels:{
		usePointStyle:true,

		generateLabels: function(chart){
			 return  chart.data.datasets.map( function( dataset, i ){
				return{
					text:dataset.label,
					lineCap:dataset.borderCapStyle,
					lineDash:[],
					lineDashOffset: 0,
					lineJoin:dataset.borderJoinStyle,
					pointStyle:'circle',
					fillStyle:dataset.backgroundColor,
					strokeStyle:dataset.borderColor,
					lineWidth:dataset.pointBorderWidth,
					lineDash:dataset.borderDash,
				}
			})
		},

	},
},

title:{
	display:true,
	text:'Востребованность сервиса',
	fontColor:'#3498db',
	fontSize:32,
	fontStyle:' bold',
	},
elements: {
	arc: {
},
	point: {radius:4,},
	line: {tension:0.4,
},
	rectangle: {
},
},
tooltips:{
},
hover:{
	mode:'nearest',
	animationDuration:400,
},
animation:{
	easing:'easeInQuart',
	},
};
 DrawTheChart(ChartData2,ChartOptions2,"chart-02","line");