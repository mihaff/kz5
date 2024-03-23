function MoreChartOptions(){} 
var ChartData3 = {
	labels : ["2021","2022","2023"],
	datasets : [
		{
		data : [389,467,588,],
		backgroundColor :'#a478f0',
		borderColor : 'rgba(136,136,136,0)',
		label:"Стартап"},

		{
		data : [214,255,279,],
		backgroundColor :'#78c0f0',
		borderColor : 'rgba(0,0,0,0)',
		label:"Копорации"},

		{
		data : [180,266,311,],
		backgroundColor :'#37e881',
		borderColor : 'rgba(156,89,182,0)',
		label:"Частные пользователи"},

]
	};
ChartOptions3= {
responsive:false,
	layout:{padding:{top:12,left:12,bottom:12,},},
	 scales: {
	xAxes:[{
		stacked: true,gridLines:{color:'rgba(0,0,0,0)',borderDash:[],},
}],

	yAxes:[{
		stacked: true,gridLines:{color:'rgba(0,0,0,0)',borderDash:[],},
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
	text:'Новые пользователи',
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
animation:{
	easing:'easeInQuart',
	},
};
 DrawTheChart(ChartData3,ChartOptions3,"chart-03","bar");