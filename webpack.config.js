const path = require('path');
const webpack = require('webpack');

module.exports = {
  devServer: {
    publicPath: '/build/',
    hot: true
  },
  entry: {
    'wildproxy': './js/wildproxy.js'
  },
  output: {
    path: path.join(__dirname, 'build'),
    filename: '[name].min.js',
  },
  resolve: {
    extensions: ['.js']
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        loaders: ['babel-loader'],
        exclude: /node_modules/
      }
    ]
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin()
  ]
};
