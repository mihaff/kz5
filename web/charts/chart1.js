function MoreChartOptions(){} 
var ChartData1 = {
	labels : ["2021","2022","2023",],
	datasets : [
		{
		data : [50,72,116,],
		backgroundColor :'#78c0f0',
		borderColor : 'rgba(136,136,136,0)',
		label:"classification"},

		{
		data : [46,65,105,],
		backgroundColor :'#37e881',
		borderColor : 'rgba(0,0,0,0)',
		label:"regression"},

]
	};
ChartOptions1= {
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
	text:'Популярность моделей',
	fontColor:'#3498db',
	fontSize:32,
	fontStyle:' bold',
	},
elements: {
	arc: {
},
	line: {
},
	rectangle: {
borderWidth:0.01,
},
},
tooltips:{
},
hover:{
	mode:'nearest',
	animationDuration:400,
},
};
 DrawTheChart(ChartData1,ChartOptions1,"chart-01","bar");