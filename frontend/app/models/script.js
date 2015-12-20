import DS from 'ember-data';

export default DS.Model.extend({
  stdio: DS.belongsTo('log'),
  status: DS.attr('string'),
  artifacts: DS.attr() // TODO: set this up
});
