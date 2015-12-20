import DS from 'ember-data';

export default DS.Model.extend({
  buildNumber: DS.attr('number'),
  branch: DS.belongsTo('branch'),
  startTime: DS.attr('date'),
  finishTime: DS.attr('date'),
  scripts: DS.hasMany('script')
});
